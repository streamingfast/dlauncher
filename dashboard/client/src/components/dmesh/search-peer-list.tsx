import React, { useState } from "react"
import { SearchOutlined } from '@ant-design/icons';
import { Peer } from "../../services/data-providers/dmesh";
import { BlockNum } from "../../atoms/block-num";
import { formatNumberWithCommas } from "../../utils/format";
import { Table, Tag } from "antd";
import Moment from "react-moment";

type Props = {
  peers: Peer[]
  headBlockNum: number
}

export const SearchPeerList: React.FC<Props> = ({ peers, headBlockNum }) => {

  const getPeers = () => {
    return peers
      .sort((a: Peer, b: Peer) => {
        if (a.tierLevel === b.tierLevel) {
          return a.host < b.host ? -1 : 1
        }
        return a.tierLevel < b.tierLevel ? -1 : 1
      })
  }


  const columns = [
    {
      title: 'Host',
      dataIndex: 'host',
      key: 'host',
    },
    {
      title: 'Tail Block',
      dataIndex: 'tail_block',
      key: 'tail_block',
      render: (text: string, peer: Peer) => (
        <BlockNum blockNum={peer.tailBlockNum} />
      ),
    },
    {
      title: 'IRR Block',
      dataIndex: 'irr_block',
      key: 'irr_block',
      render: (text: string, peer: Peer) => (
        <BlockNum blockNum={peer.irrBlockNum} />
      ),
    },
    {
      title: 'Head Block',
      dataIndex: 'head_block',
      key: 'head_block',
      render: (text: string, peer: Peer) => (
        <BlockNum blockNum={peer.headBlockNum} />
      ),
    },
    {
      title: 'Shard Size',
      dataIndex: 'shard_size',
      key: 'shard_size',
      render: (text: string, peer: Peer) => (
        formatNumberWithCommas(peer.shardSize)
      ),
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (text: string, peer: Peer) => {
        if (peer.ready) {
          return (<Tag color="green">ready</Tag>)
        } else {
            return (<Tag color="magenta">not ready</Tag>)
          }
      },
    },
    {
      title: 'Boot Time',
      dataIndex: 'boot_time',
      key: 'boot_time',
      render: (text: string, peer: Peer) => {
        if (peer.boot) {
          const date = new Date(peer.boot.seconds * 1000);
          return (
            <>
              <Moment format="YYYY-MM-DD HH:mm Z">{date.toString()}</Moment>
              <br />
              <Moment fromNow>{date.toString()}</Moment>
            </>
          )
        }
      },
    },
  ];

  return (
    <div>
      <div className="ant-table-body">
        <Table
          dataSource={getPeers()}
          columns={columns}
          pagination={{
            total: peers.length,
            pageSize: peers.length,
            hideOnSinglePage: true
          }}
        />
      </div>
    </div>
  )
}
