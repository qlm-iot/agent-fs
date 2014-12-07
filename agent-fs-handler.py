import os
import shutil
from handler import Handler
from gather_event import Event

class AgentFsHandler(Handler):
    """
    Writes events to files for agent-fs to read
    """

    def __init__(self, config):
        self.closed = False
        self.config = config
        self.stat = '/proc/stat'
        self.meminfo = '/proc/meminfo'
        self.root = config['agentfshandler.root']
        self.iis = ['user', 'nice', 'system', 'idle', 'iowait', 'irq', 'softirq']
        with open(self.stat, 'r') as f:
            for l in f:
                if l.startswith('cpu'):
                    objpath = self._make_object(l.split()[0])
                    for ii in self.iis:
                        self._make_infoitem(objpath, ii.lower())
            f.close()
        with open(self.meminfo, 'r') as f:
            objpath = self._make_object('memory')
            for l in f:
                iiname = l.split()[0][:-1]
                self._make_infoitem(objpath, iiname.lower())
            f.close()

    def _make_object(self, objname):
        objpath = self.root + '/Objects/Object.' + objname
        idpath = objpath + '/Id'
        if not os.path.exists(idpath):
            os.makedirs(idpath)
        objidfile = open(idpath + '/Text', 'w')
        objidfile.write(objname + '\n')
        objidfile.close()
        return objpath

    def _make_infoitem(self, objpath, iiname):
        iipath = objpath + '/InfoItem.' + iiname
        valpath = iipath + '/Value.0'
        if not os.path.exists(valpath):
            os.makedirs(valpath)
        iinamefile = open(iipath + '/Name', 'w')
        iinamefile.write(iiname + '\n')
        iinamefile.close()
        open(valpath + '/Text', 'a').close()
        open(valpath + '/UnixTime', 'a').close()

    def handle(self, item):
        if not self.closed:
            for member in item:
                idlist = member['id'].split('.')
                if idlist[1] == 'mem':
                    obj = 'memory'
                    ii = idlist[2]
                elif idlist[1] == 'cpu':
                    obj = idlist[2]
                    ii = idlist[3]
                iipath = self.root + '/Objects/Object.' + obj + '/InfoItem.' + ii
                ut = open(iipath + '/Value.0/UnixTime', 'w')
                ut.write(str(member['timestamp'] / 1000) + '\n')
                ut.close()
                val = open(iipath + '/Value.0/Text', 'w')
                val.write(str(member['value']) + '\n')
                val.close()

    def close(self):
        self.closed = True
        shutil.rmtree(self.root + '/Objects')
