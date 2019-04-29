import urllib.parse
import uuid
# import boto3

deal_marker = "DEAL"
table_name = "Deals-dev"

def put_items(dict_items):
    print("Found:")
    print(dict_items)
    # dynamodb = boto3.resource('dynamodb')
    # table = dynamodb.Table(table_name)
    # with table.batch_writer() as batch:
    #     for item in dict_items:
    #         item["id"] = str(uuid.uuid4())
    #         batch.put_item(
    #             Item=item
    #         )


def parse_deal(deal):
    if deal == "":
        raise ValueError("empty string")

    lines = deal.split("\n")

    # min 2: valid deal has day, description
    if len(lines) < 3:
        raise ValueError("Less than 2 lines, can't be a valid deal")
    
    deal = {}
    for line in lines:
        if line == "":
            continue
        print(line)
        # split each line to get key:value pair, 
        # only split once to prevent issues in description
        key, value = line.split(":", 1)
        deal[key] = value
    return deal


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
        # print(event)
        parsed = urllib.parse.parse_qs(request_body)
        # print("____________")
        # print(parsed)

        # TODO: how do we connect these deals to a specific location
        from_header = get_param(parsed, "from")
        print("From: %s" % from_header)

        # https://documentation.mailgun.com/en/latest/user_manual.html#routes
        # look for Parsed Messages Parameters
        subject = get_param(parsed, "subject")
        print("subject: %s\n" % subject)

        # stripped = get_param(parsed, "stripped-text")
        # print("stripped:\n%s\n" % stripped)
        # plain = get_param(parsed, "body-plain")
        # print("plain:\n%s\n" % plain)
        # html = get_param(parsed, "body-html")
        # print("html:\n%s\n" % html)

        plain = get_param(parsed, "body-plain")
        print("plain:\n%s\n" % plain)
        
        # remove all the extra \r in the string
        plain = plain.replace("\r", "")

        # split on our marker
        deals = plain.split(deal_marker)
        print(deals)
        parsed_deals = []
        for deal in deals:
            try:
                parsed_deals.append(parse_deal(deal))
            except ValueError as e:
                print("--- ValueError: %s" %e)
                pass

        if len(parsed_deals) > 0:
            put_items(parsed_deals)
        else:
            print("No valid deals found")

    except KeyError as e:
        # failed parsing
        print("!! KeyError")
        print(e)

    return {
        "statusCode": 207,
        "body": ""
    }

if __name__ == "__main__":
    event = {'resource': '/email', 'path': '/email', 'httpMethod': 'POST', 'headers': {'Accept': '*/*', 'Accept-Encoding': 'gzip', 'CloudFront-Forwarded-Proto': 'https', 'CloudFront-Is-Desktop-Viewer': 'true', 'CloudFront-Is-Mobile-Viewer': 'false', 'CloudFront-Is-SmartTV-Viewer': 'false', 'CloudFront-Is-Tablet-Viewer': 'false', 'CloudFront-Viewer-Country': 'US', 'Content-Type': 'application/x-www-form-urlencoded', 'Host': 'pvosoby5yl.execute-api.us-east-2.amazonaws.com', 'User-Agent': 'mailgun/treq-16.12.0', 'Via': '1.1 3316ddaeea3a736012726e9c08426818.cloudfront.net (CloudFront)', 'X-Amz-Cf-Id': 'DzO6tmoCLPHVvhTiKV6sFWGTznVl2lUhkLeJC5_8F6Y9HBFgb3e4vg==', 'X-Amzn-Trace-Id': 'Root=1-5cc6415e-858a2f08e7a34a3d2ea9974f', 'X-Forwarded-For': '54.144.55.146, 70.132.59.151', 'X-Forwarded-Port': '443', 'X-Forwarded-Proto': 'https'}, 'multiValueHeaders': {'Accept': ['*/*'], 'Accept-Encoding': ['gzip'], 'CloudFront-Forwarded-Proto': ['https'], 'CloudFront-Is-Desktop-Viewer': ['true'], 'CloudFront-Is-Mobile-Viewer': ['false'], 'CloudFront-Is-SmartTV-Viewer': ['false'], 'CloudFront-Is-Tablet-Viewer': ['false'], 'CloudFront-Viewer-Country': ['US'], 'Content-Type': ['application/x-www-form-urlencoded'], 'Host': ['pvosoby5yl.execute-api.us-east-2.amazonaws.com'], 'User-Agent': ['mailgun/treq-16.12.0'], 'Via': ['1.1 3316ddaeea3a736012726e9c08426818.cloudfront.net (CloudFront)'], 'X-Amz-Cf-Id': ['DzO6tmoCLPHVvhTiKV6sFWGTznVl2lUhkLeJC5_8F6Y9HBFgb3e4vg=='], 'X-Amzn-Trace-Id': ['Root=1-5cc6415e-858a2f08e7a34a3d2ea9974f'], 'X-Forwarded-For': ['54.144.55.146, 70.132.59.151'], 'X-Forwarded-Port': ['443'], 'X-Forwarded-Proto': ['https']}, 'queryStringParameters': None, 'multiValueQueryStringParameters': None, 'pathParameters': None, 'stageVariables': None, 'requestContext': {'resourceId': 'gopl5k', 'resourcePath': '/email', 'httpMethod': 'POST', 'extendedRequestId': 'Y38mwHmGiYcFr6w=', 'requestTime': '29/Apr/2019:00:12:14 +0000', 'path': '/dev/email', 'accountId': '174225498255', 'protocol': 'HTTP/1.1', 'stage': 'dev', 'domainPrefix': 'pvosoby5yl', 'requestTimeEpoch': 1556496734414, 'requestId': '70fe2f81-6a13-11e9-ae20-b7c7e6149fac', 'identity': {'cognitoIdentityPoolId': None, 'accountId': None, 'cognitoIdentityId': None, 'caller': None, 'sourceIp': '54.144.55.146', 'accessKey': None, 'cognitoAuthenticationType': None, 'cognitoAuthenticationProvider': None, 'userArn': None, 'userAgent': 'mailgun/treq-16.12.0', 'user': None}, 'domainName': 'pvosoby5yl.execute-api.us-east-2.amazonaws.com', 'apiId': 'pvosoby5yl'}, 'body': 'recipient=update%40getdealsontap.com&sender=kpquinn2%40wisc.edu&subject=this+is+a+subject&from=Kevin+Quinn+%3Ckpquinn2%40wisc.edu%3E&X-Mailgun-Incoming=Yes&X-Envelope-From=%3Ckpquinn2%40wisc.edu%3E&Received=from+wmauth3.doit.wisc.edu+%28wmauth3.doit.wisc.edu+%5B144.92.197.226%5D%29+by+mxa.mailgun.org+with+ESMTP+id+5cc6415d.7fa7a61c89f0-smtp-in-n02%3B+Mon%2C+29+Apr+2019+00%3A12%3A13+-0000+%28UTC%29&Received=from+NAM04-CO1-obe.outbound.protection.outlook.com+%28mail-co1nam04lp2058.outbound.protection.outlook.com+%5B104.47.45.58%5D%29+by+smtpauth3.wiscmail.wisc.edu+%28Oracle+Communications+Messaging+Server+8.0.1.2.20170621+64bit+%28built+Jun+21+2017%29%29+with+ESMTPS+id+%3C0PQP006W14KBPJ40%40smtpauth3.wiscmail.wisc.edu%3E+for+update%40getdealsontap.com%3B+Sun%2C+28+Apr+2019+19%3A12%3A12+-0500+%28CDT%29&X-Spam-Report=AuthenticatedSender%3Dyes%2C+SenderIP%3D%5B104.47.45.58%5D&X-Wisc-Env-From-B64=a3BxdWlubjJAd2lzYy5lZHU%3D&X-Spam-Pmxinfo=Server%3Davs-3%2C+Version%3D6.4.6.2792898%2C+Antispam-Engine%3A+2.7.2.2107409%2C+Antispam-Data%3A+2019.4.29.17%2C+AntiVirus-Engine%3A+5.60.0%2C+AntiVirus-Data%3A+2019.4.23.5600002%2C+SenderIP%3D%5B104.47.45.58%5D&Dkim-Signature=v%3D1%3B+a%3Drsa-sha256%3B+c%3Drelaxed%2Frelaxed%3B+d%3Dwisc.edu%3B+s%3Dselector1%3B+h%3DFrom%3ADate%3ASubject%3AMessage-ID%3AContent-Type%3AMIME-Version%3AX-MS-Exchange-SenderADCheck%3B+bh%3DSF56M0MUVifSSzgyex2OzppwzZjGIEoGIUGJILd2me0%3D%3B+b%3DkBoNYT51o4Bwh27HkTtc%2BS2cYcP3AApRkWdQ%2B94NoXBAgZsPxoqM%2BCR5vFCc44cu5p4KXT47yUpEPzc64ZNKqLxeFtusl7Upv8iU5AZ5oPf0WsXwW8chKIDAeqoBePgohWiaYvJXQ1GvfqV5xcuZKl%2FhmtzujzeL3zohOfCdDgc%3D&Received=from+BL0PR06MB5057.namprd06.prod.outlook.com+%2810.167.240.82%29+by+BL0PR06MB4994.namprd06.prod.outlook.com+%2810.167.235.147%29+with+Microsoft+SMTP+Server+%28version%3DTLS1_2%2C+cipher%3DTLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384%29+id+15.20.1835.12%3B+Mon%2C+29+Apr+2019+00%3A12%3A10+%2B0000&Received=from+BL0PR06MB5057.namprd06.prod.outlook.com+%28%5Bfe80%3A%3A7caa%3Adb74%3A2dc1%3A7daa%5D%29+by+BL0PR06MB5057.namprd06.prod.outlook.com+%28%5Bfe80%3A%3A7caa%3Adb74%3A2dc1%3A7daa%253%5D%29+with+mapi+id+15.20.1835.018%3B+Mon%2C+29+Apr+2019+00%3A12%3A10+%2B0000&From=Kevin+Quinn+%3Ckpquinn2%40wisc.edu%3E&To=%22update%40getdealsontap.com%22+%3Cupdate%40getdealsontap.com%3E&Subject=this+is+a+subject&Thread-Topic=this+is+a+subject&Thread-Index=AQHU%2Fh%2Bao1rMdLziKUyziAZWrvgC6g%3D%3D&Date=Mon%2C+29+Apr+2019+00%3A12%3A10+%2B0000&Message-Id=%3CBL0PR06MB50578774E1684558CBB201FCEA390%40BL0PR06MB5057.namprd06.prod.outlook.com%3E&Accept-Language=en-US&Content-Language=en-US&X-Ms-Has-Attach=&X-Ms-Tnef-Correlator=&X-Originating-Ip=%5B24.240.36.19%5D&X-Ms-Publictraffictype=Email&X-Ms-Office365-Filtering-Correlation-Id=b51253e3-a7e6-48d4-3df7-08d6cc375325&X-Ms-Office365-Filtering-Ht=Tenant&X-Microsoft-Antispam=BCL%3A0%3BPCL%3A0%3BRULEID%3A%282390118%29%287020095%29%284652040%29%288989299%29%284534185%29%284627221%29%28201703031133081%29%28201702281549075%29%288990200%29%285600141%29%28711020%29%284605104%29%284618075%29%282017052603328%29%287193020%29%3BSRVR%3ABL0PR06MB4994%3B&X-Ms-Traffictypediagnostic=BL0PR06MB4994%3A&X-Microsoft-Antispam-Prvs=%3CBL0PR06MB4994CA10D6110CDCB399DEB0EA390%40BL0PR06MB4994.namprd06.prod.outlook.com%3E&X-Ms-Oob-Tlc-Oobclassifiers=OLM%3A142%3B&X-Forefront-Prvs=0022134A87&X-Forefront-Antispam-Report=SFV%3ANSPM%3BSFS%3A%2810009020%29%28396003%29%28346002%29%28136003%29%28366004%29%28376002%29%2839860400002%29%28189003%29%28199004%29%2855016002%29%2825786009%29%2814454004%29%28558084003%29%2852536014%29%286116002%29%285660300002%29%2866066001%29%286916009%29%28256004%29%288936002%29%2881166006%29%2881156014%29%282906002%29%28105004%29%2853936002%29%2854896002%29%28486006%29%281730700003%29%28476003%29%28478600001%29%2871200400001%29%2871190400001%29%2886362001%29%283846002%29%289686003%29%2888552002%29%283480700005%29%287696005%29%288676002%29%2875432002%29%2826005%29%28102836004%29%286346003%29%28316002%29%2899286004%29%2864756008%29%2866446008%29%285640700003%29%2868736007%29%28786003%29%2897736004%29%2874316002%29%287736002%29%2833656002%29%282501003%29%2876116006%29%286506007%29%282351001%29%2819627405001%29%2873956011%29%2866946007%29%2866476007%29%28186003%29%2866556008%29%286436002%29%28212503006%29%3BDIR%3AOUT%3BSFP%3A1101%3BSCL%3A1%3BSRVR%3ABL0PR06MB4994%3BH%3ABL0PR06MB5057.namprd06.prod.outlook.com%3BFPR%3A%3BSPF%3ANone%3BLANG%3Aen%3BPTR%3AInfoNoRecords%3BMX%3A1%3BA%3A1%3B&Received-Spf=None+%28protection.outlook.com%3A+wisc.edu+does+not+designate+permitted+sender+hosts%29&Authentication-Results=spf%3Dnone+%28sender+IP+is+%29+smtp.mailfrom%3Dkpquinn2%40wisc.edu%3B&X-Ms-Exchange-Senderadcheck=1&X-Microsoft-Antispam-Message-Info=zgWmO5ZP30JGsDTKMmk2B%2BX3pFYJ2ICRFU291gmCnVGSdqxGK0xEi6YYGnLDuo46YCH%2F35EX2s8XqqOBNURyARpdSc2zTErgEGplFfA%2FphpPcmuySih6eKNUt6jnc7K%2Fd68aV24xu4%2FQVMQzIT2SXLTO0GVjBLq6AbzBuklHSq5t8alV22u5tWhdQpsF2EaGziScyBB31QekcVL%2FnLD4j4xi3lWdpbgyAr5dr4MPzvsZsGOLqczz%2FE%2FiXR7q5n0LW64LXLm%2FL4eyeNSr9zWc4CQCdmviPcKd6CdnwkMRi4197SbO55n4tLE1FmVKbJbQo29UZir5J6rTSOuI%2BzFlnVLXl4akxmhnoRNCBeGUrk6HnHFOZenbQnb8u2zWJJyzZiuySFeutjRntCem6t06vS8oPfu7GoEeCQDlSo%2FANog%3D&Content-Type=multipart%2Falternative%3B+boundary%3D%22_000_BL0PR06MB50578774E1684558CBB201FCEA390BL0PR06MB5057namp_%22&Mime-Version=1.0&X-Originatororg=wisc.edu&X-Ms-Exchange-Crosstenant-Network-Message-Id=b51253e3-a7e6-48d4-3df7-08d6cc375325&X-Ms-Exchange-Crosstenant-Originalarrivaltime=29+Apr+2019+00%3A12%3A10.8591+%28UTC%29&X-Ms-Exchange-Crosstenant-Fromentityheader=Hosted&X-Ms-Exchange-Crosstenant-Id=2ca68321-0eda-4908-88b2-424a8cb4b0f9&X-Ms-Exchange-Crosstenant-Mailboxtype=HOSTED&X-Ms-Exchange-Transport-Crosstenantheadersstamped=BL0PR06MB4994&message-headers=%5B%5B%22X-Mailgun-Incoming%22%2C+%22Yes%22%5D%2C+%5B%22X-Envelope-From%22%2C+%22%3Ckpquinn2%40wisc.edu%3E%22%5D%2C+%5B%22Received%22%2C+%22from+wmauth3.doit.wisc.edu+%28wmauth3.doit.wisc.edu+%5B144.92.197.226%5D%29+by+mxa.mailgun.org+with+ESMTP+id+5cc6415d.7fa7a61c89f0-smtp-in-n02%3B+Mon%2C+29+Apr+2019+00%3A12%3A13+-0000+%28UTC%29%22%5D%2C+%5B%22Received%22%2C+%22from+NAM04-CO1-obe.outbound.protection.outlook.com+%28mail-co1nam04lp2058.outbound.protection.outlook.com+%5B104.47.45.58%5D%29+by+smtpauth3.wiscmail.wisc.edu+%28Oracle+Communications+Messaging+Server+8.0.1.2.20170621+64bit+%28built+Jun+21+2017%29%29+with+ESMTPS+id+%3C0PQP006W14KBPJ40%40smtpauth3.wiscmail.wisc.edu%3E+for+update%40getdealsontap.com%3B+Sun%2C+28+Apr+2019+19%3A12%3A12+-0500+%28CDT%29%22%5D%2C+%5B%22X-Spam-Report%22%2C+%22AuthenticatedSender%3Dyes%2C+SenderIP%3D%5B104.47.45.58%5D%22%5D%2C+%5B%22X-Wisc-Env-From-B64%22%2C+%22a3BxdWlubjJAd2lzYy5lZHU%3D%22%5D%2C+%5B%22X-Spam-Pmxinfo%22%2C+%22Server%3Davs-3%2C+Version%3D6.4.6.2792898%2C+Antispam-Engine%3A+2.7.2.2107409%2C+Antispam-Data%3A+2019.4.29.17%2C+AntiVirus-Engine%3A+5.60.0%2C+AntiVirus-Data%3A+2019.4.23.5600002%2C+SenderIP%3D%5B104.47.45.58%5D%22%5D%2C+%5B%22Dkim-Signature%22%2C+%22v%3D1%3B+a%3Drsa-sha256%3B+c%3Drelaxed%2Frelaxed%3B+d%3Dwisc.edu%3B+s%3Dselector1%3B+h%3DFrom%3ADate%3ASubject%3AMessage-ID%3AContent-Type%3AMIME-Version%3AX-MS-Exchange-SenderADCheck%3B+bh%3DSF56M0MUVifSSzgyex2OzppwzZjGIEoGIUGJILd2me0%3D%3B+b%3DkBoNYT51o4Bwh27HkTtc%2BS2cYcP3AApRkWdQ%2B94NoXBAgZsPxoqM%2BCR5vFCc44cu5p4KXT47yUpEPzc64ZNKqLxeFtusl7Upv8iU5AZ5oPf0WsXwW8chKIDAeqoBePgohWiaYvJXQ1GvfqV5xcuZKl%2FhmtzujzeL3zohOfCdDgc%3D%22%5D%2C+%5B%22Received%22%2C+%22from+BL0PR06MB5057.namprd06.prod.outlook.com+%2810.167.240.82%29+by+BL0PR06MB4994.namprd06.prod.outlook.com+%2810.167.235.147%29+with+Microsoft+SMTP+Server+%28version%3DTLS1_2%2C+cipher%3DTLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384%29+id+15.20.1835.12%3B+Mon%2C+29+Apr+2019+00%3A12%3A10+%2B0000%22%5D%2C+%5B%22Received%22%2C+%22from+BL0PR06MB5057.namprd06.prod.outlook.com+%28%5Bfe80%3A%3A7caa%3Adb74%3A2dc1%3A7daa%5D%29+by+BL0PR06MB5057.namprd06.prod.outlook.com+%28%5Bfe80%3A%3A7caa%3Adb74%3A2dc1%3A7daa%253%5D%29+with+mapi+id+15.20.1835.018%3B+Mon%2C+29+Apr+2019+00%3A12%3A10+%2B0000%22%5D%2C+%5B%22From%22%2C+%22Kevin+Quinn+%3Ckpquinn2%40wisc.edu%3E%22%5D%2C+%5B%22To%22%2C+%22%5C%22update%40getdealsontap.com%5C%22+%3Cupdate%40getdealsontap.com%3E%22%5D%2C+%5B%22Subject%22%2C+%22this+is+a+subject%22%5D%2C+%5B%22Thread-Topic%22%2C+%22this+is+a+subject%22%5D%2C+%5B%22Thread-Index%22%2C+%22AQHU%2Fh%2Bao1rMdLziKUyziAZWrvgC6g%3D%3D%22%5D%2C+%5B%22Date%22%2C+%22Mon%2C+29+Apr+2019+00%3A12%3A10+%2B0000%22%5D%2C+%5B%22Message-Id%22%2C+%22%3CBL0PR06MB50578774E1684558CBB201FCEA390%40BL0PR06MB5057.namprd06.prod.outlook.com%3E%22%5D%2C+%5B%22Accept-Language%22%2C+%22en-US%22%5D%2C+%5B%22Content-Language%22%2C+%22en-US%22%5D%2C+%5B%22X-Ms-Has-Attach%22%2C+%22%22%5D%2C+%5B%22X-Ms-Tnef-Correlator%22%2C+%22%22%5D%2C+%5B%22X-Originating-Ip%22%2C+%22%5B24.240.36.19%5D%22%5D%2C+%5B%22X-Ms-Publictraffictype%22%2C+%22Email%22%5D%2C+%5B%22X-Ms-Office365-Filtering-Correlation-Id%22%2C+%22b51253e3-a7e6-48d4-3df7-08d6cc375325%22%5D%2C+%5B%22X-Ms-Office365-Filtering-Ht%22%2C+%22Tenant%22%5D%2C+%5B%22X-Microsoft-Antispam%22%2C+%22BCL%3A0%3BPCL%3A0%3BRULEID%3A%282390118%29%287020095%29%284652040%29%288989299%29%284534185%29%284627221%29%28201703031133081%29%28201702281549075%29%288990200%29%285600141%29%28711020%29%284605104%29%284618075%29%282017052603328%29%287193020%29%3BSRVR%3ABL0PR06MB4994%3B%22%5D%2C+%5B%22X-Ms-Traffictypediagnostic%22%2C+%22BL0PR06MB4994%3A%22%5D%2C+%5B%22X-Microsoft-Antispam-Prvs%22%2C+%22%3CBL0PR06MB4994CA10D6110CDCB399DEB0EA390%40BL0PR06MB4994.namprd06.prod.outlook.com%3E%22%5D%2C+%5B%22X-Ms-Oob-Tlc-Oobclassifiers%22%2C+%22OLM%3A142%3B%22%5D%2C+%5B%22X-Forefront-Prvs%22%2C+%220022134A87%22%5D%2C+%5B%22X-Forefront-Antispam-Report%22%2C+%22SFV%3ANSPM%3BSFS%3A%2810009020%29%28396003%29%28346002%29%28136003%29%28366004%29%28376002%29%2839860400002%29%28189003%29%28199004%29%2855016002%29%2825786009%29%2814454004%29%28558084003%29%2852536014%29%286116002%29%285660300002%29%2866066001%29%286916009%29%28256004%29%288936002%29%2881166006%29%2881156014%29%282906002%29%28105004%29%2853936002%29%2854896002%29%28486006%29%281730700003%29%28476003%29%28478600001%29%2871200400001%29%2871190400001%29%2886362001%29%283846002%29%289686003%29%2888552002%29%283480700005%29%287696005%29%288676002%29%2875432002%29%2826005%29%28102836004%29%286346003%29%28316002%29%2899286004%29%2864756008%29%2866446008%29%285640700003%29%2868736007%29%28786003%29%2897736004%29%2874316002%29%287736002%29%2833656002%29%282501003%29%2876116006%29%286506007%29%282351001%29%2819627405001%29%2873956011%29%2866946007%29%2866476007%29%28186003%29%2866556008%29%286436002%29%28212503006%29%3BDIR%3AOUT%3BSFP%3A1101%3BSCL%3A1%3BSRVR%3ABL0PR06MB4994%3BH%3ABL0PR06MB5057.namprd06.prod.outlook.com%3BFPR%3A%3BSPF%3ANone%3BLANG%3Aen%3BPTR%3AInfoNoRecords%3BMX%3A1%3BA%3A1%3B%22%5D%2C+%5B%22Received-Spf%22%2C+%22None+%28protection.outlook.com%3A+wisc.edu+does+not+designate+permitted+sender+hosts%29%22%5D%2C+%5B%22Authentication-Results%22%2C+%22spf%3Dnone+%28sender+IP+is+%29+smtp.mailfrom%3Dkpquinn2%40wisc.edu%3B%22%5D%2C+%5B%22X-Ms-Exchange-Senderadcheck%22%2C+%221%22%5D%2C+%5B%22X-Microsoft-Antispam-Message-Info%22%2C+%22zgWmO5ZP30JGsDTKMmk2B%2BX3pFYJ2ICRFU291gmCnVGSdqxGK0xEi6YYGnLDuo46YCH%2F35EX2s8XqqOBNURyARpdSc2zTErgEGplFfA%2FphpPcmuySih6eKNUt6jnc7K%2Fd68aV24xu4%2FQVMQzIT2SXLTO0GVjBLq6AbzBuklHSq5t8alV22u5tWhdQpsF2EaGziScyBB31QekcVL%2FnLD4j4xi3lWdpbgyAr5dr4MPzvsZsGOLqczz%2FE%2FiXR7q5n0LW64LXLm%2FL4eyeNSr9zWc4CQCdmviPcKd6CdnwkMRi4197SbO55n4tLE1FmVKbJbQo29UZir5J6rTSOuI%2BzFlnVLXl4akxmhnoRNCBeGUrk6HnHFOZenbQnb8u2zWJJyzZiuySFeutjRntCem6t06vS8oPfu7GoEeCQDlSo%2FANog%3D%22%5D%2C+%5B%22Content-Type%22%2C+%22multipart%2Falternative%3B+boundary%3D%5C%22_000_BL0PR06MB50578774E1684558CBB201FCEA390BL0PR06MB5057namp_%5C%22%22%5D%2C+%5B%22Mime-Version%22%2C+%221.0%22%5D%2C+%5B%22X-Originatororg%22%2C+%22wisc.edu%22%5D%2C+%5B%22X-Ms-Exchange-Crosstenant-Network-Message-Id%22%2C+%22b51253e3-a7e6-48d4-3df7-08d6cc375325%22%5D%2C+%5B%22X-Ms-Exchange-Crosstenant-Originalarrivaltime%22%2C+%2229+Apr+2019+00%3A12%3A10.8591+%28UTC%29%22%5D%2C+%5B%22X-Ms-Exchange-Crosstenant-Fromentityheader%22%2C+%22Hosted%22%5D%2C+%5B%22X-Ms-Exchange-Crosstenant-Id%22%2C+%222ca68321-0eda-4908-88b2-424a8cb4b0f9%22%5D%2C+%5B%22X-Ms-Exchange-Crosstenant-Mailboxtype%22%2C+%22HOSTED%22%5D%2C+%5B%22X-Ms-Exchange-Transport-Crosstenantheadersstamped%22%2C+%22BL0PR06MB4994%22%5D%5D&timestamp=1556496733&token=cbd8b808ec7613d6b9c49f72bb1b17e26feda6d4a6cd888cb6&signature=55f035a8261b84c04ba624bc747a9488a9d28f054649d500c95d8ad6ae3bc08b&body-plain=DEAL%0D%0Adays%3AM%2CW%2CF%0D%0Atime%3A8-12%0D%0Adesc%3A+hella+deals%0D%0A%0D%0ADEAL%0D%0Adays%3AM%2CW%2CF%0D%0Atime%3A8-12%0D%0Adesc%3A+sick+deal+bro%0D%0A%0D%0A&body-html=%3Chtml%3E%0D%0A%3Chead%3E%0D%0A%3Cmeta+http-equiv%3D%22Content-Type%22+content%3D%22text%2Fhtml%3B+charset%3Diso-8859-1%22%3E%0D%0A%3Cstyle+type%3D%22text%2Fcss%22+style%3D%22display%3Anone%3B%22%3E+P+%7Bmargin-top%3A0%3Bmargin-bottom%3A0%3B%7D+%3C%2Fstyle%3E%0D%0A%3C%2Fhead%3E%0D%0A%3Cbody+dir%3D%22ltr%22%3E%0D%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%3B+color%3A+rgb%280%2C+0%2C+0%29%3B%22%3E%0D%0ADEAL%3C%2Fdiv%3E%0D%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%3B+color%3A+rgb%280%2C+0%2C+0%29%3B%22%3E%0D%0Adays%3AM%2CW%2CF%3C%2Fdiv%3E%0D%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%3B+color%3A+rgb%280%2C+0%2C+0%29%3B%22%3E%0D%0Atime%3A8-12%3C%2Fdiv%3E%0D%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%3B+color%3A+rgb%280%2C+0%2C+0%29%3B%22%3E%0D%0Adesc%3A+hella+deals%3C%2Fdiv%3E%0D%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%3B+color%3A+rgb%280%2C+0%2C+0%29%3B%22%3E%0D%0A%3Cbr%3E%0D%0A%3C%2Fdiv%3E%0D%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%3B+color%3A+rgb%280%2C+0%2C+0%29%3B%22%3E%0D%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%22%3E%0D%0ADEAL%3C%2Fdiv%3E%0D%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%22%3E%0D%0Adays%3AM%2CW%2CF%3C%2Fdiv%3E%0D%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%22%3E%0D%0Atime%3A8-12%3C%2Fdiv%3E%0D%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%22%3E%0D%0Adesc%3A+sick+deal+bro%3Cbr%3E%0D%0A%3C%2Fdiv%3E%0D%0A%3Cbr%3E%0D%0A%3C%2Fdiv%3E%0D%0A%3C%2Fbody%3E%0D%0A%3C%2Fhtml%3E%0D%0A&stripped-html=%3Chtml%3E%0A%3Chead%3E%0A%3Cmeta+http-equiv%3D%22Content-Type%22+content%3D%22text%2Fhtml%3B+charset%3Diso-8859-1%22%3E%0A%3Cstyle+type%3D%22text%2Fcss%22+style%3D%22display%3Anone%3B%22%3E+P+%7Bmargin-top%3A0%3Bmargin-bottom%3A0%3B%7D+%3C%2Fstyle%3E%0A%3C%2Fhead%3E%0A%3Cbody+dir%3D%22ltr%22%3E%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%3B+color%3A+rgb%280%2C+0%2C+0%29%3B%22%3E%0ADEAL%3C%2Fdiv%3E%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%3B+color%3A+rgb%280%2C+0%2C+0%29%3B%22%3E%0Adays%3AM%2CW%2CF%3C%2Fdiv%3E%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%3B+color%3A+rgb%280%2C+0%2C+0%29%3B%22%3E%0Atime%3A8-12%3C%2Fdiv%3E%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%3B+color%3A+rgb%280%2C+0%2C+0%29%3B%22%3E%0Adesc%3A+hella+deals%3C%2Fdiv%3E%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%3B+color%3A+rgb%280%2C+0%2C+0%29%3B%22%3E%0A%3Cbr%3E%0A%3C%2Fdiv%3E%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%3B+color%3A+rgb%280%2C+0%2C+0%29%3B%22%3E%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%22%3E%0ADEAL%3C%2Fdiv%3E%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%22%3E%0Adays%3AM%2CW%2CF%3C%2Fdiv%3E%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%22%3E%0Atime%3A8-12%3C%2Fdiv%3E%0A%3Cdiv+style%3D%22font-family%3A+Calibri%2C+Arial%2C+Helvetica%2C+sans-serif%3B+font-size%3A+12pt%22%3E%0Adesc%3A+sick+deal+bro%3Cbr%3E%0A%3C%2Fdiv%3E%0A%3Cbr%3E%0A%3C%2Fdiv%3E%0A%3C%2Fbody%3E%0A%3C%2Fhtml%3E%0A&stripped-text=DEAL%0D%0Adays%3AM%2CW%2CF%0D%0Atime%3A8-12%0D%0Adesc%3A+hella+deals%0D%0A%0D%0ADEAL%0D%0Adays%3AM%2CW%2CF%0D%0Atime%3A8-12%0D%0Adesc%3A+sick+deal+bro&stripped-signature=', 'isBase64Encoded': False}
    # parsed_body = {'recipient': ['update@getdealsontap.com'], 'sender': ['kpquinn2@wisc.edu'], 'subject': ['this is a subject'], 'from': ['Kevin Quinn <kpquinn2@wisc.edu>'], 'X-Mailgun-Incoming': ['Yes'], 'X-Envelope-From': ['<kpquinn2@wisc.edu>'], 'Received': ['from wmauth3.doit.wisc.edu (wmauth3.doit.wisc.edu [144.92.197.226]) by mxa.mailgun.org with ESMTP id 5cc6415d.7fa7a61c89f0-smtp-in-n02; Mon, 29 Apr 2019 00:12:13 -0000 (UTC)', 'from NAM04-CO1-obe.outbound.protection.outlook.com (mail-co1nam04lp2058.outbound.protection.outlook.com [104.47.45.58]) by smtpauth3.wiscmail.wisc.edu (Oracle Communications Messaging Server 8.0.1.2.20170621 64bit (built Jun 21 2017)) with ESMTPS id <0PQP006W14KBPJ40@smtpauth3.wiscmail.wisc.edu> for update@getdealsontap.com; Sun, 28 Apr 2019 19:12:12 -0500 (CDT)', 'from BL0PR06MB5057.namprd06.prod.outlook.com (10.167.240.82) by BL0PR06MB4994.namprd06.prod.outlook.com (10.167.235.147) with Microsoft SMTP Server (version=TLS1_2, cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id 15.20.1835.12; Mon, 29 Apr 2019 00:12:10 +0000', 'from BL0PR06MB5057.namprd06.prod.outlook.com ([fe80::7caa:db74:2dc1:7daa]) by BL0PR06MB5057.namprd06.prod.outlook.com ([fe80::7caa:db74:2dc1:7daa%3]) with mapi id 15.20.1835.018; Mon, 29 Apr 2019 00:12:10 +0000'], 'X-Spam-Report': ['AuthenticatedSender=yes, SenderIP=[104.47.45.58]'], 'X-Wisc-Env-From-B64': ['a3BxdWlubjJAd2lzYy5lZHU='], 'X-Spam-Pmxinfo': ['Server=avs-3, Version=6.4.6.2792898, Antispam-Engine: 2.7.2.2107409, Antispam-Data: 2019.4.29.17, AntiVirus-Engine: 5.60.0, AntiVirus-Data: 2019.4.23.5600002, SenderIP=[104.47.45.58]'], 'Dkim-Signature': ['v=1; a=rsa-sha256; c=relaxed/relaxed; d=wisc.edu; s=selector1; h=From:Date:Subject:Message-ID:Content-Type:MIME-Version:X-MS-Exchange-SenderADCheck; bh=SF56M0MUVifSSzgyex2OzppwzZjGIEoGIUGJILd2me0=; b=kBoNYT51o4Bwh27HkTtc+S2cYcP3AApRkWdQ+94NoXBAgZsPxoqM+CR5vFCc44cu5p4KXT47yUpEPzc64ZNKqLxeFtusl7Upv8iU5AZ5oPf0WsXwW8chKIDAeqoBePgohWiaYvJXQ1GvfqV5xcuZKl/hmtzujzeL3zohOfCdDgc='], 'From': ['Kevin Quinn <kpquinn2@wisc.edu>'], 'To': ['"update@getdealsontap.com" <update@getdealsontap.com>'], 'Subject': ['this is a subject'], 'Thread-Topic': ['this is a subject'], 'Thread-Index': ['AQHU/h+ao1rMdLziKUyziAZWrvgC6g=='], 'Date': ['Mon, 29 Apr 2019 00:12:10 +0000'], 'Message-Id': ['<BL0PR06MB50578774E1684558CBB201FCEA390@BL0PR06MB5057.namprd06.prod.outlook.com>'], 'Accept-Language': ['en-US'], 'Content-Language': ['en-US'], 'X-Originating-Ip': ['[24.240.36.19]'], 'X-Ms-Publictraffictype': ['Email'], 'X-Ms-Office365-Filtering-Correlation-Id': ['b51253e3-a7e6-48d4-3df7-08d6cc375325'], 'X-Ms-Office365-Filtering-Ht': ['Tenant'], 'X-Microsoft-Antispam': ['BCL:0;PCL:0;RULEID:(2390118)(7020095)(4652040)(8989299)(4534185)(4627221)(201703031133081)(201702281549075)(8990200)(5600141)(711020)(4605104)(4618075)(2017052603328)(7193020);SRVR:BL0PR06MB4994;'], 'X-Ms-Traffictypediagnostic': ['BL0PR06MB4994:'], 'X-Microsoft-Antispam-Prvs': ['<BL0PR06MB4994CA10D6110CDCB399DEB0EA390@BL0PR06MB4994.namprd06.prod.outlook.com>'], 'X-Ms-Oob-Tlc-Oobclassifiers': ['OLM:142;'], 'X-Forefront-Prvs': ['0022134A87'], 'X-Forefront-Antispam-Report': ['SFV:NSPM;SFS:(10009020)(396003)(346002)(136003)(366004)(376002)(39860400002)(189003)(199004)(55016002)(25786009)(14454004)(558084003)(52536014)(6116002)(5660300002)(66066001)(6916009)(256004)(8936002)(81166006)(81156014)(2906002)(105004)(53936002)(54896002)(486006)(1730700003)(476003)(478600001)(71200400001)(71190400001)(86362001)(3846002)(9686003)(88552002)(3480700005)(7696005)(8676002)(75432002)(26005)(102836004)(6346003)(316002)(99286004)(64756008)(66446008)(5640700003)(68736007)(786003)(97736004)(74316002)(7736002)(33656002)(2501003)(76116006)(6506007)(2351001)(19627405001)(73956011)(66946007)(66476007)(186003)(66556008)(6436002)(212503006);DIR:OUT;SFP:1101;SCL:1;SRVR:BL0PR06MB4994;H:BL0PR06MB5057.namprd06.prod.outlook.com;FPR:;SPF:None;LANG:en;PTR:InfoNoRecords;MX:1;A:1;'], 'Received-Spf': ['None (protection.outlook.com: wisc.edu does not designate permitted sender hosts)'], 'Authentication-Results': ['spf=none (sender IP is ) smtp.mailfrom=kpquinn2@wisc.edu;'], 'X-Ms-Exchange-Senderadcheck': ['1'], 'X-Microsoft-Antispam-Message-Info': ['zgWmO5ZP30JGsDTKMmk2B+X3pFYJ2ICRFU291gmCnVGSdqxGK0xEi6YYGnLDuo46YCH/35EX2s8XqqOBNURyARpdSc2zTErgEGplFfA/phpPcmuySih6eKNUt6jnc7K/d68aV24xu4/QVMQzIT2SXLTO0GVjBLq6AbzBuklHSq5t8alV22u5tWhdQpsF2EaGziScyBB31QekcVL/nLD4j4xi3lWdpbgyAr5dr4MPzvsZsGOLqczz/E/iXR7q5n0LW64LXLm/L4eyeNSr9zWc4CQCdmviPcKd6CdnwkMRi4197SbO55n4tLE1FmVKbJbQo29UZir5J6rTSOuI+zFlnVLXl4akxmhnoRNCBeGUrk6HnHFOZenbQnb8u2zWJJyzZiuySFeutjRntCem6t06vS8oPfu7GoEeCQDlSo/ANog='], 'Content-Type': ['multipart/alternative; boundary="_000_BL0PR06MB50578774E1684558CBB201FCEA390BL0PR06MB5057namp_"'], 'Mime-Version': ['1.0'], 'X-Originatororg': ['wisc.edu'], 'X-Ms-Exchange-Crosstenant-Network-Message-Id': ['b51253e3-a7e6-48d4-3df7-08d6cc375325'], 'X-Ms-Exchange-Crosstenant-Originalarrivaltime': ['29 Apr 2019 00:12:10.8591 (UTC)'], 'X-Ms-Exchange-Crosstenant-Fromentityheader': ['Hosted'], 'X-Ms-Exchange-Crosstenant-Id': ['2ca68321-0eda-4908-88b2-424a8cb4b0f9'], 'X-Ms-Exchange-Crosstenant-Mailboxtype': ['HOSTED'], 'X-Ms-Exchange-Transport-Crosstenantheadersstamped': ['BL0PR06MB4994'], 'message-headers': ['[["X-Mailgun-Incoming", "Yes"], ["X-Envelope-From", "<kpquinn2@wisc.edu>"], ["Received", "from wmauth3.doit.wisc.edu (wmauth3.doit.wisc.edu [144.92.197.226]) by mxa.mailgun.org with ESMTP id 5cc6415d.7fa7a61c89f0-smtp-in-n02; Mon, 29 Apr 2019 00:12:13 -0000 (UTC)"], ["Received", "from NAM04-CO1-obe.outbound.protection.outlook.com (mail-co1nam04lp2058.outbound.protection.outlook.com [104.47.45.58]) by smtpauth3.wiscmail.wisc.edu (Oracle Communications Messaging Server 8.0.1.2.20170621 64bit (built Jun 21 2017)) with ESMTPS id <0PQP006W14KBPJ40@smtpauth3.wiscmail.wisc.edu> for update@getdealsontap.com; Sun, 28 Apr 2019 19:12:12 -0500 (CDT)"], ["X-Spam-Report", "AuthenticatedSender=yes, SenderIP=[104.47.45.58]"], ["X-Wisc-Env-From-B64", "a3BxdWlubjJAd2lzYy5lZHU="], ["X-Spam-Pmxinfo", "Server=avs-3, Version=6.4.6.2792898, Antispam-Engine: 2.7.2.2107409, Antispam-Data: 2019.4.29.17, AntiVirus-Engine: 5.60.0, AntiVirus-Data: 2019.4.23.5600002, SenderIP=[104.47.45.58]"], ["Dkim-Signature", "v=1; a=rsa-sha256; c=relaxed/relaxed; d=wisc.edu; s=selector1; h=From:Date:Subject:Message-ID:Content-Type:MIME-Version:X-MS-Exchange-SenderADCheck; bh=SF56M0MUVifSSzgyex2OzppwzZjGIEoGIUGJILd2me0=; b=kBoNYT51o4Bwh27HkTtc+S2cYcP3AApRkWdQ+94NoXBAgZsPxoqM+CR5vFCc44cu5p4KXT47yUpEPzc64ZNKqLxeFtusl7Upv8iU5AZ5oPf0WsXwW8chKIDAeqoBePgohWiaYvJXQ1GvfqV5xcuZKl/hmtzujzeL3zohOfCdDgc="], ["Received", "from BL0PR06MB5057.namprd06.prod.outlook.com (10.167.240.82) by BL0PR06MB4994.namprd06.prod.outlook.com (10.167.235.147) with Microsoft SMTP Server (version=TLS1_2, cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id 15.20.1835.12; Mon, 29 Apr 2019 00:12:10 +0000"], ["Received", "from BL0PR06MB5057.namprd06.prod.outlook.com ([fe80::7caa:db74:2dc1:7daa]) by BL0PR06MB5057.namprd06.prod.outlook.com ([fe80::7caa:db74:2dc1:7daa%3]) with mapi id 15.20.1835.018; Mon, 29 Apr 2019 00:12:10 +0000"], ["From", "Kevin Quinn <kpquinn2@wisc.edu>"], ["To", "\\"update@getdealsontap.com\\" <update@getdealsontap.com>"], ["Subject", "this is a subject"], ["Thread-Topic", "this is a subject"], ["Thread-Index", "AQHU/h+ao1rMdLziKUyziAZWrvgC6g=="], ["Date", "Mon, 29 Apr 2019 00:12:10 +0000"], ["Message-Id", "<BL0PR06MB50578774E1684558CBB201FCEA390@BL0PR06MB5057.namprd06.prod.outlook.com>"], ["Accept-Language", "en-US"], ["Content-Language", "en-US"], ["X-Ms-Has-Attach", ""], ["X-Ms-Tnef-Correlator", ""], ["X-Originating-Ip", "[24.240.36.19]"], ["X-Ms-Publictraffictype", "Email"], ["X-Ms-Office365-Filtering-Correlation-Id", "b51253e3-a7e6-48d4-3df7-08d6cc375325"], ["X-Ms-Office365-Filtering-Ht", "Tenant"], ["X-Microsoft-Antispam", "BCL:0;PCL:0;RULEID:(2390118)(7020095)(4652040)(8989299)(4534185)(4627221)(201703031133081)(201702281549075)(8990200)(5600141)(711020)(4605104)(4618075)(2017052603328)(7193020);SRVR:BL0PR06MB4994;"], ["X-Ms-Traffictypediagnostic", "BL0PR06MB4994:"], ["X-Microsoft-Antispam-Prvs", "<BL0PR06MB4994CA10D6110CDCB399DEB0EA390@BL0PR06MB4994.namprd06.prod.outlook.com>"], ["X-Ms-Oob-Tlc-Oobclassifiers", "OLM:142;"], ["X-Forefront-Prvs", "0022134A87"], ["X-Forefront-Antispam-Report", "SFV:NSPM;SFS:(10009020)(396003)(346002)(136003)(366004)(376002)(39860400002)(189003)(199004)(55016002)(25786009)(14454004)(558084003)(52536014)(6116002)(5660300002)(66066001)(6916009)(256004)(8936002)(81166006)(81156014)(2906002)(105004)(53936002)(54896002)(486006)(1730700003)(476003)(478600001)(71200400001)(71190400001)(86362001)(3846002)(9686003)(88552002)(3480700005)(7696005)(8676002)(75432002)(26005)(102836004)(6346003)(316002)(99286004)(64756008)(66446008)(5640700003)(68736007)(786003)(97736004)(74316002)(7736002)(33656002)(2501003)(76116006)(6506007)(2351001)(19627405001)(73956011)(66946007)(66476007)(186003)(66556008)(6436002)(212503006);DIR:OUT;SFP:1101;SCL:1;SRVR:BL0PR06MB4994;H:BL0PR06MB5057.namprd06.prod.outlook.com;FPR:;SPF:None;LANG:en;PTR:InfoNoRecords;MX:1;A:1;"], ["Received-Spf", "None (protection.outlook.com: wisc.edu does not designate permitted sender hosts)"], ["Authentication-Results", "spf=none (sender IP is ) smtp.mailfrom=kpquinn2@wisc.edu;"], ["X-Ms-Exchange-Senderadcheck", "1"], ["X-Microsoft-Antispam-Message-Info", "zgWmO5ZP30JGsDTKMmk2B+X3pFYJ2ICRFU291gmCnVGSdqxGK0xEi6YYGnLDuo46YCH/35EX2s8XqqOBNURyARpdSc2zTErgEGplFfA/phpPcmuySih6eKNUt6jnc7K/d68aV24xu4/QVMQzIT2SXLTO0GVjBLq6AbzBuklHSq5t8alV22u5tWhdQpsF2EaGziScyBB31QekcVL/nLD4j4xi3lWdpbgyAr5dr4MPzvsZsGOLqczz/E/iXR7q5n0LW64LXLm/L4eyeNSr9zWc4CQCdmviPcKd6CdnwkMRi4197SbO55n4tLE1FmVKbJbQo29UZir5J6rTSOuI+zFlnVLXl4akxmhnoRNCBeGUrk6HnHFOZenbQnb8u2zWJJyzZiuySFeutjRntCem6t06vS8oPfu7GoEeCQDlSo/ANog="], ["Content-Type", "multipart/alternative; boundary=\\"_000_BL0PR06MB50578774E1684558CBB201FCEA390BL0PR06MB5057namp_\\""], ["Mime-Version", "1.0"], ["X-Originatororg", "wisc.edu"], ["X-Ms-Exchange-Crosstenant-Network-Message-Id", "b51253e3-a7e6-48d4-3df7-08d6cc375325"], ["X-Ms-Exchange-Crosstenant-Originalarrivaltime", "29 Apr 2019 00:12:10.8591 (UTC)"], ["X-Ms-Exchange-Crosstenant-Fromentityheader", "Hosted"], ["X-Ms-Exchange-Crosstenant-Id", "2ca68321-0eda-4908-88b2-424a8cb4b0f9"], ["X-Ms-Exchange-Crosstenant-Mailboxtype", "HOSTED"], ["X-Ms-Exchange-Transport-Crosstenantheadersstamped", "BL0PR06MB4994"]]'], 'timestamp': ['1556496733'], 'token': ['cbd8b808ec7613d6b9c49f72bb1b17e26feda6d4a6cd888cb6'], 'signature': ['55f035a8261b84c04ba624bc747a9488a9d28f054649d500c95d8ad6ae3bc08b'], 'body-plain': ['DEAL\r\ndays:M,W,F\r\ntime:8-12\r\ndesc: hella deals\r\n\r\nDEAL\r\ndays:M,W,F\r\ntime:8-12\r\ndesc: sick deal bro\r\n\r\n'], 'body-html': ['<html>\r\n<head>\r\n<meta http-equiv="Content-Type" content="text/html; charset=iso-8859-1">\r\n<style type="text/css" style="display:none;"> P {margin-top:0;margin-bottom:0;} </style>\r\n</head>\r\n<body dir="ltr">\r\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt; color: rgb(0, 0, 0);">\r\nDEAL</div>\r\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt; color: rgb(0, 0, 0);">\r\ndays:M,W,F</div>\r\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt; color: rgb(0, 0, 0);">\r\ntime:8-12</div>\r\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt; color: rgb(0, 0, 0);">\r\ndesc: hella deals</div>\r\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt; color: rgb(0, 0, 0);">\r\n<br>\r\n</div>\r\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt; color: rgb(0, 0, 0);">\r\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt">\r\nDEAL</div>\r\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt">\r\ndays:M,W,F</div>\r\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt">\r\ntime:8-12</div>\r\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt">\r\ndesc: sick deal bro<br>\r\n</div>\r\n<br>\r\n</div>\r\n</body>\r\n</html>\r\n'], 'stripped-html': ['<html>\n<head>\n<meta http-equiv="Content-Type" content="text/html; charset=iso-8859-1">\n<style type="text/css" style="display:none;"> P {margin-top:0;margin-bottom:0;} </style>\n</head>\n<body dir="ltr">\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt; color: rgb(0, 0, 0);">\nDEAL</div>\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt; color: rgb(0, 0, 0);">\ndays:M,W,F</div>\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt; color: rgb(0, 0, 0);">\ntime:8-12</div>\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt; color: rgb(0, 0, 0);">\ndesc: hella deals</div>\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt; color: rgb(0, 0, 0);">\n<br>\n</div>\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt; color: rgb(0, 0, 0);">\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt">\nDEAL</div>\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt">\ndays:M,W,F</div>\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt">\ntime:8-12</div>\n<div style="font-family: Calibri, Arial, Helvetica, sans-serif; font-size: 12pt">\ndesc: sick deal bro<br>\n</div>\n<br>\n</div>\n</body>\n</html>\n'], 'stripped-text': ['DEAL\r\ndays:M,W,F\r\ntime:8-12\r\ndesc: hella deals\r\n\r\nDEAL\r\ndays:M,W,F\r\ntime:8-12\r\ndesc: sick deal bro']}
    process(event, None)