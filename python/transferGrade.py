
from pylti.pylti import common
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
    response_url = "https://courses.edx.org/courses/course-v1:MITx+15.071x+1T2019/xblock/block-v1:MITx+15.071x+1T2019+type@lti_consumer+block@a855518774854399b79abee373351e3c/handler_noauth/outcome_service_handler"
    lis_result_sourcedid = "course-v1%3AMITx%2B15.071x%2B1T2019:courses.edx.org-a855518774854399b79abee373351e3c:6987787dd79cf0aecabdca8ddae95b4a"
    is_success = post_grade(0.5, lis_result_sourcedid, response_url)
    print(is_success)

if __name__ == "__main__":
    main()
