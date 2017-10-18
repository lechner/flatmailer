package main

//     Copyright (C) 2017 Felix Lechner <felix.lechner@lease-up.com>
//     This file is part of Flatmailer.

//     Flatmailer is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 2 of the License, or
//     (at your option) any later version.

//     Flatmailer is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.

//     You should have received a copy of the GNU General Public License
//     along with Flatmailer.  If not, see <http://www.gnu.org/licenses/>.

import "flag"
import "fmt"
import "strconv"
import "log"
import "bufio"
import "path"
import "os"
import "time"
import "syscall"

var diagnosticCode, envelopeId, remoteServer, quoteLimit string
var firstAttempt, mostRecentAttempt, failAfter string

func init() {

	// string flags
	flag.StringVar(&diagnosticCode, "diagnostic-code", "",
		"Diagnostic code `MESSAGE`")
	flag.StringVar(&envelopeId, "envelope-id", "",
		"Original envelope `ID`")
	flag.StringVar(&remoteServer, "remote-server", "",
		"`HOST` name of remote server")

	// int64 flags
	flag.StringVar(&firstAttempt, "first-attempt", "",
		"`TIMESTAMP` of first send attempt (default: filesystem mtime)")
	flag.StringVar(&mostRecentAttempt, "most-recent-attempt", "",
		"`TIMESTAMP` of most recent send attempt (default: filesystem atime)")
	flag.StringVar(&failAfter, "fail-after", "",
		"`TIMESTAMP` after which to declare failure")

	// int flags
	flag.StringVar(&quoteLimit, "quote-limit", "",
		"Most `LINES` included from original message (default: all)")
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] status-code < message\n",
			path.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  Produces a delivery status notification (DSN) for a queued message.\n")
		fmt.Fprintf(os.Stderr, "  The status should look like 4.#.# for a delay, or 5.#.# for a failure.\n")
		fmt.Fprintf(os.Stderr, "  All timestamps should be in RFC3339 format.\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "The flags are:\n")
		flag.PrintDefaults()
	}

}

// TODO
// utf-8 for addressing and message type may need to be declared

func main() {

	// configure logger
	log.SetPrefix("flatmailer-dsn: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// parse command line flags
	flag.Parse()

	ourInputFormat := time.RFC3339
	ourOutputFormat := time.RFC1123Z
	
	if len(flag.Args()) != 1 {
		log.Fatalln("Need exactly one argument besides flags")
	}	
	status := flag.Arg(0)

	var class, subject, detail int
	items, err := fmt.Sscanf(status, "%1d.%3d.%3d",
		&class, &subject, &detail)
	if items != 3 || (class != 4 && class != 5) {
		log.Fatalln("Status must have format 4.#.# or 5.#.#:", err);
	}

	delay := (class == 4)

	// stat input file
	info, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalln("Could not stat source message:", err)
		// original: exit(111)
	}

	// set time of first attempt
	if firstAttempt != "" {
		// rewrite timestamp in our format
		timestamp, err := time.Parse(ourInputFormat, firstAttempt)
		if err != nil {
			log.Fatalln("timestamp", firstAttempt,
				"has invalid format:", err)
		}
		firstAttempt = timestamp.Format(ourOutputFormat)
	} else {
		// get mtime
		firstAttempt = info.ModTime().Format(ourOutputFormat)
		// original: ctime
	}

	if mostRecentAttempt != "" {
		timestamp, err := time.Parse(ourInputFormat, mostRecentAttempt)
		if err != nil {
			log.Fatalln("timestamp", mostRecentAttempt,
				"has invalid format:", err)
		}
		mostRecentAttempt = timestamp.Format(ourOutputFormat)
	} else {
		// get atime
		stat := info.Sys().(*syscall.Stat_t)
		timestamp := time.Unix(stat.Atim.Sec, stat.Atim.Nsec)
		mostRecentAttempt = timestamp.Format(ourOutputFormat)
	}

	if failAfter != "" {
		timestamp, err := time.Parse(ourInputFormat, failAfter);
		if err != nil {
			log.Fatalln("timestamp", failAfter,
				"has invalid format:", err)
		}
		failAfter = timestamp.Format(ourOutputFormat)
	}
	
	if quoteLimit == "" {
		// config_readint("bouncelines", opt_lines);
	}
	maxLines := 0
	if quoteLimit != "" {
		maxLines, err = strconv.Atoi(quoteLimit)
		if err != nil {
			log.Fatalln("parameter for quote limit must me a number")
		}
	}

	doubleBounceTo := "felix.lechner@lease-up.com"
	//   if (!config_read("doublebounceto", doublebounceto)
	//       || !doublebounceto)
	//     config_read("adminaddr", doublebounceto);

	mailHostId := "dummy"
	//   read_hostnames();
	//   if (!config_read("idhost", idhost))
	//     idhost = me;
	//   else
	//     canonicalize(idhost);

	bounceTo := ""
	//   config_read("bounceto", bounceto);

	// get time for timestamps
	thisTime := time.Now()
	boundary := fmt.Sprintf("%d.%9d.%d", thisTime.Unix(),
		thisTime.Nanosecond(), os.Getpid())
	// original: microseconds with six digits
	messageId := fmt.Sprintf("<%s.flatmailer@%s>", boundary, mailHostId)
		
	// read from stdin
	scanner := bufio.NewScanner(os.Stdin)

	// read original sender
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatalln("Could not read sender address from message:",
				err)
		}
	}
	sender:= scanner.Text()	

	// read original recipients
	var recipients []string
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		recipients = append(recipients, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln("No recipients were read from message", err)
	}

	// find best recipient for bounce message
	if sender != "" {
		if bounceTo != "" {
			fmt.Printf("%s", bounceTo)
		} else {
			fmt.Printf("%s", sender)
		}
	} else {
		if doubleBounceTo != "" {
			fmt.Printf("#@[]\n")
			fmt.Printf("%s", doubleBounceTo)
		} else {
			log.Fatalln("Nowhere to send double bounce")
		}
	}
	fmt.Printf("\n")

	// generate delivery status notification
	
	fmt.Printf("\n")
	fmt.Printf("From: Message Delivery Subsystem <MAILER-DAEMON@%s\n",
		mailHostId)
	fmt.Printf("To: <%s>\n", sender)
	fmt.Printf("Subject: Returned mail: Could not send message\n")
	fmt.Printf("Date: %s\n", thisTime.Format(ourOutputFormat))
	fmt.Printf("Message-Id: %s\n", messageId)

	fmt.Printf("MIME-Version: 1.0\n")
	fmt.Printf("Content-Type: multipart/report; report-type=delivery-status;\n")
	fmt.Printf("\tboundary=\"%s\"\n", boundary)

	// human readable text portion

	fmt.Printf("\n")
	fmt.Printf("--%s\n", boundary)
	fmt.Printf("Content-Type: text/plain; charset=us-ascii\n")
	fmt.Printf("\n")
	fmt.Printf("This is the flatmailer delivery system. The message attached below \n")
	if delay {
		fmt.Printf("has not been")
	} else {
		fmt.Printf("could not be")
	}
	fmt.Printf(" delivered to one or more of the intended recipients:\n")
	
	fmt.Printf("\n")
	for _, recipient := range recipients {
		fmt.Printf("\t<%s>\n", recipient)
	}
	
	if delay {
		if failAfter != "" {
			fmt.Printf("\n")
			fmt.Printf("Delivery attempts will continue until %s\n",
				failAfter)
		}
		fmt.Printf("\n")
		fmt.Printf("A final delivery status notification will be generated if delivery\n")
		fmt.Printf("proves to be impossible within the configured time limit.\n")
	}
	
	// delivery-status portion

	fmt.Printf("\n")
	fmt.Printf("--%s\n", boundary)
	fmt.Printf("Content-Type: message/delivery-status\n")
	
	fmt.Printf("\n")
	fmt.Printf("Reporting-MTA: x-local-hostname; %s\n", mailHostId)
	fmt.Printf("Arrival-Date: %s\n", firstAttempt)
	
	if envelopeId != "" {
		fmt.Printf("Original-Envelope-Id: %s\n", envelopeId)
	}

	for _, recipient := range recipients {
		fmt.Printf("\n")
		fmt.Printf("Final-Recipient: rfc822; %s\n", recipient)
		if delay {
			fmt.Printf("Action: delayed\n")
		} else {
			fmt.Printf("Action: failed\n")
		}
		fmt.Printf("Status: %s\n", status)
		fmt.Printf("Last-Attempt-Date: %s\n", mostRecentAttempt)

		if remoteServer != "" {
			fmt.Printf("Remote-MTA: dns; %s\n", remoteServer)
		}
		if diagnosticCode != "" {
			fmt.Printf("Diagnostic-Code: %s\n", diagnosticCode)
		}
		if delay && failAfter != "" {
			fmt.Printf("Will-Retry-Until: %s\n", failAfter)
		}
	}
	
	// copy the original message
	fmt.Printf("\n")
	fmt.Printf("--%s\n", boundary)
	fmt.Printf("Content-Type: message/rfc822")
	fmt.Printf("\n")
	
	// copy header
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}		
		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln("Could not read message header:", err)
	}

	// optionally copy body
	if quoteLimit == "" || maxLines > 0 {
		fmt.Printf("\n")
		lines := 0
		for (quoteLimit == "" || lines < maxLines) && scanner.Scan() {
			fmt.Println(scanner.Text())
			lines++
		}
		if err := scanner.Err(); err != nil {
			log.Fatalln("Could not read message body:", err)
		}
	}
	
	fmt.Printf("\n")
	fmt.Printf("--%s\n", boundary)
}
