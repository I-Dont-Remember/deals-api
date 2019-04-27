import urllib.parse

def get_param(body, name):
    try:
        # query params come as lists
        return body[name][0]
    except KeyError as e:
        print("[!] KeyError %s" %e)
        return None

def process(event, context):
    # process email from mailgun
    try:
        request_body = event["body"]
        print(event)
        parsed = urllib.parse.parse_qs(request_body)
        print("____________")
        print(parsed)

        from_header = get_param(parsed, "from")
        print("From: %s" % from_header)

        # https://documentation.mailgun.com/en/latest/user_manual.html#routes
        # look for Parsed Messages Parameters
        subject = get_param(parsed, "subject")
        print("subject: %s\n" % subject)

        stripped = get_param(parsed, "stripped-text")
        print("stripped:\n%s\n" % stripped)

        plain = get_param(parsed, "body-plain")
        print("plain:\n%s\n" % plain)

        # could split text first based on a marker, such as DEAL
        # then split by newline to get each key:value pair
        # then split by colon to get each part
        # ! need to remove all the \r in the strings
        html = get_param(parsed, "body-html")
        print("html:\n%s\n" % html)

    except KeyError as e:
        # failed parsing
        print("!! KeyError")
        print(e)

    return {
        "statusCode": 207,
        "body": ""
    }

if __name__ == "__main__":
    event = {'resource': '/email', 'path': '/email', 'httpMethod': 'POST', 'headers': {'Accept': '*/*', 'Accept-Encoding': 'gzip', 'CloudFront-Forwarded-Proto': 'https', 'CloudFront-Is-Desktop-Viewer': 'true', 'CloudFront-Is-Mobile-Viewer': 'false', 'CloudFront-Is-SmartTV-Viewer': 'false', 'CloudFront-Is-Tablet-Viewer': 'false', 'CloudFront-Viewer-Country': 'US', 'Content-Type': 'application/x-www-form-urlencoded', 'Host': 'pvosoby5yl.execute-api.us-east-2.amazonaws.com', 'User-Agent': 'mailgun/treq-16.12.0', 'Via': '1.1 de6760156d781e28f72545a2e9243b26.cloudfront.net (CloudFront)', 'X-Amz-Cf-Id': 'mqBv_ZmzFXYpxFgvo6slU4SsgYAlw6owsRqvLLmQTqeg8CDGM-1t9w==', 'X-Amzn-Trace-Id': 'Root=1-5cc4bf5a-e9f26db377068483addc3ec7', 'X-Forwarded-For': '100.26.117.121, 70.132.32.91', 'X-Forwarded-Port': '443', 'X-Forwarded-Proto': 'https'}, 'multiValueHeaders': {'Accept': ['*/*'], 'Accept-Encoding': ['gzip'], 'CloudFront-Forwarded-Proto': ['https'], 'CloudFront-Is-Desktop-Viewer': ['true'], 'CloudFront-Is-Mobile-Viewer': ['false'], 'CloudFront-Is-SmartTV-Viewer': ['false'], 'CloudFront-Is-Tablet-Viewer': ['false'], 'CloudFront-Viewer-Country': ['US'], 'Content-Type': ['application/x-www-form-urlencoded'], 'Host': ['pvosoby5yl.execute-api.us-east-2.amazonaws.com'], 'User-Agent': ['mailgun/treq-16.12.0'], 'Via': ['1.1 de6760156d781e28f72545a2e9243b26.cloudfront.net (CloudFront)'], 'X-Amz-Cf-Id': ['mqBv_ZmzFXYpxFgvo6slU4SsgYAlw6owsRqvLLmQTqeg8CDGM-1t9w=='], 'X-Amzn-Trace-Id': ['Root=1-5cc4bf5a-e9f26db377068483addc3ec7'], 'X-Forwarded-For': ['100.26.117.121, 70.132.32.91'], 'X-Forwarded-Port': ['443'], 'X-Forwarded-Proto': ['https']}, 'queryStringParameters': None, 'multiValueQueryStringParameters': None, 'pathParameters': None, 'stageVariables': None, 'requestContext': {'resourceId': 'gopl5k', 'resourcePath': '/email', 'httpMethod': 'POST', 'extendedRequestId': 'Y0LWJGkTCYcFvSQ=', 'requestTime': '27/Apr/2019:20:45:14 +0000', 'path': '/dev/email', 'accountId': '174225498255', 'protocol': 'HTTP/1.1', 'stage': 'dev', 'domainPrefix': 'pvosoby5yl', 'requestTimeEpoch': 1556397914554, 'requestId': '5bc4327c-692d-11e9-ac59-45df291c739b', 'identity': {'cognitoIdentityPoolId': None, 'accountId': None, 'cognitoIdentityId': None, 'caller': None, 'sourceIp': '100.26.117.121', 'accessKey': None, 'cognitoAuthenticationType': None, 'cognitoAuthenticationProvider': None, 'userArn': None, 'userAgent': 'mailgun/treq-16.12.0', 'user': None}, 'domainName': 'pvosoby5yl.execute-api.us-east-2.amazonaws.com', 'apiId': 'pvosoby5yl'}, 'body': 'recipient=update%40getdealsontap.com&sender=quinn618%40gmail.com&subject=testing+lambda&from=Kevin+Quinn+%3Cquinn618%40gmail.com%3E&X-Mailgun-Incoming=Yes&X-Envelope-From=%3Cquinn618%40gmail.com%3E&Received=from+mail-ed1-f46.google.com+%28mail-ed1-f46.google.com+%5B209.85.208.46%5D%29+by+mxa.mailgun.org+with+ESMTP+id+5cc4bf57.7f14dc0517f0-smtp-in-n02%3B+Sat%2C+27+Apr+2019+20%3A45%3A11+-0000+%28UTC%29&Received=by+mail-ed1-f46.google.com+with+SMTP+id+a6so6013840edv.1++++++++for+%3Cupdate%40getdealsontap.com%3E%3B+Sat%2C+27+Apr+2019+13%3A45%3A11+-0700+%28PDT%29&Dkim-Signature=v%3D1%3B+a%3Drsa-sha256%3B+c%3Drelaxed%2Frelaxed%3B++++++++d%3Dgmail.com%3B+s%3D20161025%3B++++++++h%3Dmime-version%3Afrom%3Adate%3Amessage-id%3Asubject%3Ato%3B++++++++bh%3DJhnQGWZgdS9j7I%2FqCJP8Wk2Y7SoIx4RhtZU%2BjsAP6sM%3D%3B++++++++b%3DaAYGy9pL4A1C%2F3ARA2EvD7nxJVN8HU9OSbZSTw3OFDqR1H4jK4qsudf%2FDnsc%2B26rw%2B+++++++++7uPIr6Vj28rNYsz%2BOT69SNlGkZ15LRjIE2lU4YyR7FWskoetOcph5DV3dvk6N6enCWOB+++++++++mXgGZhYznBm5hUQH%2BrxFJKf5n5%2BaTSR2Ua5QkYlghiEF8AfKsrMqQSsWuvAscCwa7VIo+++++++++CfOTBEtu8TkMjKNox8XkBENmh0rET03Q0xDYHq5akAbJm4xJMhRwQYaqJDgTvHT%2BdT1H+++++++++hL2ggZdDtotFZowF2ZYaKKpY1HeU7%2BHhzOqid3ugteuuAOjWl8ZMlTgqy9tAIrf3AFxq+++++++++pidQ%3D%3D&X-Google-Dkim-Signature=v%3D1%3B+a%3Drsa-sha256%3B+c%3Drelaxed%2Frelaxed%3B++++++++d%3D1e100.net%3B+s%3D20161025%3B++++++++h%3Dx-gm-message-state%3Amime-version%3Afrom%3Adate%3Amessage-id%3Asubject%3Ato%3B++++++++bh%3DJhnQGWZgdS9j7I%2FqCJP8Wk2Y7SoIx4RhtZU%2BjsAP6sM%3D%3B++++++++b%3DC%2BJVPbPBvOy2DkLNG%2B3mcw5xCM72rUl3wFPnAkaj%2FFmgrsBpRpIqKr1iXFFtkbrn0n+++++++++%2Ff%2FJV7o4PvCjfn1Kstc1fyyJ2N%2BUxjlWIXxs9%2BfgJmqRTR0w5LsAbkB4gHm5hfH7BDGB+++++++++Q4YFmueCTu6b4Zr8EWKHIeUZb7N6X3yq6IuYJfw%2FQqHBbnBkIlAYU9Y6UGJZ32pxh4CS+++++++++itCcjwTA%2FTnRK9mvr9AeouxwHFEd4mZsRovLW7ywttlJOiowtYhXCIMfX%2BIQN9sZmtVW+++++++++mb%2B2ujXPon5CygU3HAR5O%2BmWKg1dgR7fXtVu0kjbh4boIotIlS5FM%2F6Xak8w708PCa8a+++++++++Q2Ng%3D%3D&X-Gm-Message-State=APjAAAU8YLbZf4a11vbQkDAeFyI1DJT81aDcrSuyyJXWOPzhkU%2Bh5cdF%09cFFstYuBfR7ygvRjvycTu2Dvf7xkdFlryLYsBG3iZDk%2F&X-Google-Smtp-Source=APXvYqyEcYv9O1x6KDtQkesCrERa0JT0HgPb1A%2FdS9vWF%2FoQOqDZU68KmisioWdi2xLAPFgr3MaQObSAaTOQUCF95vQ%3D&X-Received=by+2002%3Aa50%3A996d%3A%3A+with+SMTP+id+l42mr33381880edb.181.1556397910661%3B+Sat%2C+27+Apr+2019+13%3A45%3A10+-0700+%28PDT%29&Mime-Version=1.0&From=Kevin+Quinn+%3Cquinn618%40gmail.com%3E&Date=Sat%2C+27+Apr+2019+15%3A44%3A59+-0500&Message-Id=%3CCALU2L_MrABeDWJzho2gO6eBXQkGRsDs6Ni7vpP-OPYzSRxOSxg%40mail.gmail.com%3E&Subject=testing+lambda&To=update%40getdealsontap.com&Content-Type=multipart%2Falternative%3B+boundary%3D%2200000000000098213a0587892017%22&message-headers=%5B%5B%22X-Mailgun-Incoming%22%2C+%22Yes%22%5D%2C+%5B%22X-Envelope-From%22%2C+%22%3Cquinn618%40gmail.com%3E%22%5D%2C+%5B%22Received%22%2C+%22from+mail-ed1-f46.google.com+%28mail-ed1-f46.google.com+%5B209.85.208.46%5D%29+by+mxa.mailgun.org+with+ESMTP+id+5cc4bf57.7f14dc0517f0-smtp-in-n02%3B+Sat%2C+27+Apr+2019+20%3A45%3A11+-0000+%28UTC%29%22%5D%2C+%5B%22Received%22%2C+%22by+mail-ed1-f46.google.com+with+SMTP+id+a6so6013840edv.1++++++++for+%3Cupdate%40getdealsontap.com%3E%3B+Sat%2C+27+Apr+2019+13%3A45%3A11+-0700+%28PDT%29%22%5D%2C+%5B%22Dkim-Signature%22%2C+%22v%3D1%3B+a%3Drsa-sha256%3B+c%3Drelaxed%2Frelaxed%3B++++++++d%3Dgmail.com%3B+s%3D20161025%3B++++++++h%3Dmime-version%3Afrom%3Adate%3Amessage-id%3Asubject%3Ato%3B++++++++bh%3DJhnQGWZgdS9j7I%2FqCJP8Wk2Y7SoIx4RhtZU%2BjsAP6sM%3D%3B++++++++b%3DaAYGy9pL4A1C%2F3ARA2EvD7nxJVN8HU9OSbZSTw3OFDqR1H4jK4qsudf%2FDnsc%2B26rw%2B+++++++++7uPIr6Vj28rNYsz%2BOT69SNlGkZ15LRjIE2lU4YyR7FWskoetOcph5DV3dvk6N6enCWOB+++++++++mXgGZhYznBm5hUQH%2BrxFJKf5n5%2BaTSR2Ua5QkYlghiEF8AfKsrMqQSsWuvAscCwa7VIo+++++++++CfOTBEtu8TkMjKNox8XkBENmh0rET03Q0xDYHq5akAbJm4xJMhRwQYaqJDgTvHT%2BdT1H+++++++++hL2ggZdDtotFZowF2ZYaKKpY1HeU7%2BHhzOqid3ugteuuAOjWl8ZMlTgqy9tAIrf3AFxq+++++++++pidQ%3D%3D%22%5D%2C+%5B%22X-Google-Dkim-Signature%22%2C+%22v%3D1%3B+a%3Drsa-sha256%3B+c%3Drelaxed%2Frelaxed%3B++++++++d%3D1e100.net%3B+s%3D20161025%3B++++++++h%3Dx-gm-message-state%3Amime-version%3Afrom%3Adate%3Amessage-id%3Asubject%3Ato%3B++++++++bh%3DJhnQGWZgdS9j7I%2FqCJP8Wk2Y7SoIx4RhtZU%2BjsAP6sM%3D%3B++++++++b%3DC%2BJVPbPBvOy2DkLNG%2B3mcw5xCM72rUl3wFPnAkaj%2FFmgrsBpRpIqKr1iXFFtkbrn0n+++++++++%2Ff%2FJV7o4PvCjfn1Kstc1fyyJ2N%2BUxjlWIXxs9%2BfgJmqRTR0w5LsAbkB4gHm5hfH7BDGB+++++++++Q4YFmueCTu6b4Zr8EWKHIeUZb7N6X3yq6IuYJfw%2FQqHBbnBkIlAYU9Y6UGJZ32pxh4CS+++++++++itCcjwTA%2FTnRK9mvr9AeouxwHFEd4mZsRovLW7ywttlJOiowtYhXCIMfX%2BIQN9sZmtVW+++++++++mb%2B2ujXPon5CygU3HAR5O%2BmWKg1dgR7fXtVu0kjbh4boIotIlS5FM%2F6Xak8w708PCa8a+++++++++Q2Ng%3D%3D%22%5D%2C+%5B%22X-Gm-Message-State%22%2C+%22APjAAAU8YLbZf4a11vbQkDAeFyI1DJT81aDcrSuyyJXWOPzhkU%2Bh5cdF%5CtcFFstYuBfR7ygvRjvycTu2Dvf7xkdFlryLYsBG3iZDk%2F%22%5D%2C+%5B%22X-Google-Smtp-Source%22%2C+%22APXvYqyEcYv9O1x6KDtQkesCrERa0JT0HgPb1A%2FdS9vWF%2FoQOqDZU68KmisioWdi2xLAPFgr3MaQObSAaTOQUCF95vQ%3D%22%5D%2C+%5B%22X-Received%22%2C+%22by+2002%3Aa50%3A996d%3A%3A+with+SMTP+id+l42mr33381880edb.181.1556397910661%3B+Sat%2C+27+Apr+2019+13%3A45%3A10+-0700+%28PDT%29%22%5D%2C+%5B%22Mime-Version%22%2C+%221.0%22%5D%2C+%5B%22From%22%2C+%22Kevin+Quinn+%3Cquinn618%40gmail.com%3E%22%5D%2C+%5B%22Date%22%2C+%22Sat%2C+27+Apr+2019+15%3A44%3A59+-0500%22%5D%2C+%5B%22Message-Id%22%2C+%22%3CCALU2L_MrABeDWJzho2gO6eBXQkGRsDs6Ni7vpP-OPYzSRxOSxg%40mail.gmail.com%3E%22%5D%2C+%5B%22Subject%22%2C+%22testing+lambda%22%5D%2C+%5B%22To%22%2C+%22update%40getdealsontap.com%22%5D%2C+%5B%22Content-Type%22%2C+%22multipart%2Falternative%3B+boundary%3D%5C%2200000000000098213a0587892017%5C%22%22%5D%5D&timestamp=1556397911&token=7ce17c643a5e50bc1a1d1ca747e0f766e2a2384fe761735233&signature=cf2d53411c3807849649f47ba21032b7c6602687c50180b810c337ea626cb9a5&body-plain=lambda+test+boi%0D%0A&body-html=%3Cdiv+dir%3D%22ltr%22%3Elambda+test+boi%3Cbr%3E%3C%2Fdiv%3E%0D%0A&stripped-html=%3Cdiv+dir%3D%22ltr%22%3Elambda+test+boi%3Cbr%3E%3C%2Fdiv%3E%0A&stripped-text=lambda+test+boi&stripped-signature=', 'isBase64Encoded': False}
    process(event, None)