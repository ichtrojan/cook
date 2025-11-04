package mailer

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/textproto"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hibiken/asynq"
	"github.com/ichtrojan/cook/config"
)

type EmailPayload struct {
	TemplateName string
	To           string
	Subject      string
	Data         map[string]interface{}
	Attachments  []*Attachment
}

type Attachment struct {
	Filename    string
	ContentType string
	Content     []byte
}

func Send(ctx context.Context, t *asynq.Task) (err error) {
	var payload EmailPayload

	if err = json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %v", err)
	}

	tmpl, err := template.ParseFiles(fmt.Sprintf("./templates/emails/%s.html", payload.TemplateName))

	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	buf := new(bytes.Buffer)

	if err = tmpl.Execute(buf, payload.Data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	htmlBody := buf.String()

	var msgBuffer bytes.Buffer

	writer := multipart.NewWriter(&msgBuffer)

	headers := make(map[string][]string)
	headers["From"] = []string{fmt.Sprintf("%s <%s>", config.AppConfig.AppName, config.MailConfig.FromAddress)}
	headers["To"] = []string{payload.To}
	headers["Subject"] = []string{payload.Subject}
	headers["Content-Type"] = []string{"multipart/mixed; boundary=" + writer.Boundary()}

	for key, value := range headers {
		for _, v := range value {
			msgBuffer.WriteString(fmt.Sprintf("%s: %s\r\n", key, v))
		}
	}

	msgBuffer.WriteString("\r\n")

	bodyPart, _ := writer.CreatePart(textproto.MIMEHeader{"Content-Type": []string{"text/html; charset=UTF-8"}})

	_, err = bodyPart.Write([]byte(htmlBody))

	if err != nil {
		return err
	}

	for _, attachment := range payload.Attachments {
		if attachment == nil {
			continue
		}

		encodedAttachment := base64.StdEncoding.EncodeToString(attachment.Content)

		attachmentHeader := textproto.MIMEHeader{
			"Content-Type":              []string{attachment.ContentType},
			"Content-Transfer-Encoding": []string{"base64"},
			"Content-Disposition":       []string{fmt.Sprintf("attachment; filename=\"%s\"", attachment.Filename)},
		}

		attachmentPart, _ := writer.CreatePart(attachmentHeader)

		_, err = attachmentPart.Write([]byte(encodedAttachment))

		if err != nil {
			return err
		}
	}

	err = writer.Close()

	if err != nil {
		return err
	}

	rawEmail := &ses.RawMessage{
		Data: msgBuffer.Bytes(),
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AwsConfig.Region),
		Credentials: credentials.NewStaticCredentials(config.AwsConfig.AccessKeyID, config.AwsConfig.SecretAccessKey, ""),
	})

	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	svc := ses.New(sess)

	input := &ses.SendRawEmailInput{
		RawMessage: rawEmail,
	}

	_, err = svc.SendRawEmail(input)

	if err != nil {
		return err
	}

	return nil
}

func EnqueueEmailTask(client *asynq.Client, payload EmailPayload) error {
	data, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	task := asynq.NewTask("send:email", data)

	_, err = client.Enqueue(task)

	if err != nil {
		return fmt.Errorf("failed to enqueue task: %v", err)
	}

	return nil
}
