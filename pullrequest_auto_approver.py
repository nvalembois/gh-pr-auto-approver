import logging
import os
import re
import warnings
import time
from pprint import pprint
from github import Github

warnings.filterwarnings("ignore")

logging.basicConfig(
    format='%(levelname) -5s %(asctime)s : %(message)s',
    datefmt='%d-%b-%y %H:%M:%S',
    level=logging.INFO)

log = logging.getLogger(__name__) 

github_repo = os.environ.get('GITHUB_REPO')
if not re.match("\\w+/\\w+", github_repo):
    log.error("Environment Variables GITHUB_REPO not passed correctly")
    raise SystemExit(1)
    
github_token = os.environ.get('GITHUB_TOKEN')
if not re.match("\\w+", github_token):
    log.error("Environment Variables GITHUB_TOKEN not passed correctly")
    raise SystemExit(1)

github_base = os.environ.get('GITHUB_BASE', 'main')
if not re.match("\\w+", github_base):
    log.error("Environment Variables GITHUB_BASE not passed correctly")
    raise SystemExit(1)

if __name__ == "__main__":
    log.info('PR auto approver : start')

    log.info('Connect Github')
    g = Github(github_token)
    
    log.info('get repo')
    repo = g.get_repo(github_repo)

    log.info('List openned PR for branch main')
    pulls = repo.get_pulls(state='open', sort='created', base=github_base)

    for idx, pr in enumerate(pulls):
        log.info('Merge PR #{num}: {title}'.format(num=pr.number, title=pr.title))
        # if not pr.mergeable:
        #     log.info('! PR not mergeable: {reason}'.format(reason=pr.mergeable_state))
        #     continue
        res = pr.merge()
        if not res.merged:
            log.warn('PR failed to merge: {reason}'.format(reason=res.message))

    log.info('PR auto approver : end')