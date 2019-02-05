import networkx as nx
import re

for level in range(1,5):
    for asmt in ["cc","rts","dfe"]:
        print(level,asmt)
        lnr = nx.nx_pydot.read_dot('./python/model/{}/{}/tree.dot'.format(level,asmt))
        old_label = nx.get_node_attributes(lnr,'label')
        new_label = {}
        for key,value in old_label.items():
            res = re.search('<B> (.*)</B>',value)
            if res is None:
                res = re.search('<B>(.*)</B>',value)
            text = res.group(1)
            text = ' '.join(text.split())
            text = text.strip()
            new_label[key] = text
        nx.set_node_attributes(lnr, new_label, 'label')
        nx.nx_pydot.write_dot(lnr, './python/model/{}/{}/pytree.dot'.format(level,asmt))
