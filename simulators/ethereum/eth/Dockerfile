ARG branch=latest
FROM ethereum/client-go:$branch

ADD genesis.json /genesis.json
ADD chain.rlp.gz /chain.rlp.gz

RUN /usr/local/bin/geth init /genesis.json
RUN /usr/local/bin/geth import /chain.rlp.gz

ENTRYPOINT /usr/local/bin/geth --verbosity=5 --nat=none --nodiscover --nousb