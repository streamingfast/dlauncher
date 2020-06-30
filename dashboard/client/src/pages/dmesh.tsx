import React, { useEffect, useState } from "react"
import { Row, Col, Typography, PageHeader, Descriptions, Tag, Button } from "antd"
import { SearchPeerList } from "../components/dmesh/search-peer-list";
import { getDmesh, Peer } from "../services/data-providers/dmesh";
import { DisconnectOutlined, SyncOutlined } from "@ant-design/icons/lib";
import { withBaseLayout } from "../components/layout/layout";
import { retryFunc } from "../utils/retry";
import { useInterval } from "../hooks/use-interval";

const { Text } = Typography

export const BaseDmeshPage: React.FC = () => {
  const [connected, setConnected] = useState(false)
  const [peers, setPeers] = useState<Peer[]>([])
  const [headBlockNum, setHeadBlockNum] = useState(0)


  const tryGetDmesh = async () => {
    const res = await getDmesh();
    if (!res || !res.clientsList || res.clientsList.length <= 0)
      throw new Error('peer list empty');
    const peers = res.clientsList
      .map(peer => {
        return {
          ...peer,
          // statusString: AppStatusNumberToStringMap[app.status]
        };
      })
    setPeers(peers);
  };


  useInterval(() => {
    setConnected(true)
    tryGetDmesh()
  }, 1000);


  return (
    <Row gutter={[16, 16]}>
      <Col span={24}>
        <PageHeader
          title="Dmesh Peers"
          tags={
            connected ? (
              <Tag color="geekblue">
                <SyncOutlined spin /> connected
              </Tag>
            ) : (
              <Tag color="red">
                <DisconnectOutlined /> disconnected
              </Tag>
            )
          }
        >
          <Descriptions size="small" column={3}>
            <Descriptions.Item label="Watch Key">
              <Text code>
                /local
              </Text>
            </Descriptions.Item>
          </Descriptions>
        </PageHeader>
        <Row>
          <Col>
            <div style={{ marginTop: "10px" }}>
              <SearchPeerList peers={peers} headBlockNum={headBlockNum} />
            </div>
          </Col>
        </Row>
      </Col>
    </Row>
  )
}
export const DmeshPage = withBaseLayout(BaseDmeshPage);
