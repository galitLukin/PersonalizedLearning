#import pylti
#from pylti import common
from ../pylti.pylti import common
from pylti.pylti.common import LTIPostMessageException

def post_grade(grade, lis_result_sourcedid, response_url):
    """
    Post grade to LTI consumer using XML
    :param: grade: 0 <= grade <= 1
    :return: True if post successful and grade valid
    :exception: LTIPostMessageException if call failed
    """

    _consumers = {
        "oandgkey": {
            "secret": "oandgsecret",
            "cert": None
            }
    }
    message_identifier_id = "edX_fix"
    operation = 'replaceResult'
    # # edX devbox fix
    score = float(grade)
    if 0 <= score <= 1.0:
        xml = common.generate_request_xml(message_identifier_id, operation, lis_result_sourcedid, score)
        ret = common.post_message(_consumers, "oandgkey", response_url, xml)
        if not ret:
            raise LTIPostMessageException("Post Message Failed")
        return True
    return False



def main():
    response_url = state = sys.argv[1]
    lis_result_sourcedid = sys.argv[2]
    is_success = post_grade(0.5, lis_result_sourcedid, response_url)
    print(is_success)

if __name__ == "__main__":
    main()
