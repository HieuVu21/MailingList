package main

import (
	"context"
	"log"
	pb "mailinglist/Proto"
	"time"

	"github.com/alexflint/go-arg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func LogResponse(res *pb.EmailResponse, err error) {
	if err != nil {
		log.Fatal("err")
	}
	if res.EmailEntry == nil {
		log.Print("email not found ")
	} else {
		log.Printf("response: %v", res.EmailEntry)
	}
}
func CreateEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry {
	log.Printf("create Email: %v", addr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.CreateEmail(ctx, &pb.CreateEmailRequest{EmailAddr: addr})
	LogResponse(res, err)
	return res.EmailEntry
}

func GetEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry  {
	log.Printf("get Email: %v",addr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	res, err := client.GetEmail(ctx, &pb.GetEmailRequest{EmailAddr: addr})
	LogResponse(res, err)
	return res.EmailEntry
}
func GetEmailBatch(client pb.MailingListServiceClient, count int, page int) {
	log.Printf("create Email batch: %v %v ",count,page)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetEmailBatch(ctx, &pb.GetEmailBatchRequest{Count: int32(count), Page: int32(page)})
	if err != nil {
		log.Fatal("err")
	}
	log.Println("respond")
	for i := 0; i < len(res.EmailEntry); i++ {
		log.Print()
	}
}

func UpdateEmail(client pb.MailingListServiceClient, entry pb.EmailEntry) *pb.EmailEntry {
	log.Printf("update Email:" )
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.UpdateEmail(ctx, &pb.UpdateEmailRequest{EmailEntry: &entry})
	LogResponse(res, err)
	return res.EmailEntry
	
}
func DeleteEmai( client pb.MailingListServiceClient, addr string) *pb.EmailEntry {
	log.Printf("delete Email")
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()

	res , err := client.DeleteEmail(ctx, &pb.DeleteEmailRequest{EmailAddr: addr})
	LogResponse(res,err)
	return res.EmailEntry
}


var agrs struct{
	GrpcAddr string `arg:env:MAILINGLIST_GRPC_ADDR`
}
func main(){
	arg.MustParse(&agrs)
	if agrs.GrpcAddr ==""{
		agrs.GrpcAddr = ":8080"
	}
	conn ,err := grpc.Dial(agrs.GrpcAddr,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		log.Fatal("not connect")
	}
	defer conn.Close()
	client := pb.NewMailingListServiceClient(conn)

	newEmail := CreateEmail(client,"hieuvu3@gmail.com")
	newEmail.ConfirmedAt = 10000
	UpdateEmail(client, *newEmail)
	DeleteEmai(client,newEmail.Email )
	// GetEmailBatch(client,1,5)
}