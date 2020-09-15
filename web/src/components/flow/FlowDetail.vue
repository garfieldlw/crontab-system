<template>
    <div>
        <h3 align="left">工作流 > 详情</h3>
        <div align="left">
            <a-descriptions layout="vertical" :column="4" bordered>
                <a-descriptions-item label="Id">{{data.id}}</a-descriptions-item>
                <a-descriptions-item label="名称">{{data.name}}</a-descriptions-item>
                <a-descriptions-item label="类型（1周期任务，2即时任务）">{{data.flow_type}}</a-descriptions-item>
                <a-descriptions-item label="计划时间">{{data.spec}}</a-descriptions-item>
                <a-descriptions-item label="状态(1可用, 2禁用)">{{data.status}}</a-descriptions-item>
                <a-descriptions-item label="描述">{{data.desc}}</a-descriptions-item>
                <a-descriptions-item label="创建时间">{{DateTimeFormat(data.create_time)}}</a-descriptions-item>
                <a-descriptions-item label="更新时间">{{DateTimeFormat(data.update_time)}}</a-descriptions-item>
                <a-descriptions-item label="工作流信息" :span="4">{{data.info}}</a-descriptions-item>
                <a-descriptions-item label="日志" :span="4">
                    <a-table
                            :columns="columns"
                            :data-source="dataLogFlow"
                            :rowKey="record => record.id.toString()"
                            :pagination="pagination"
                            :loading="loading"
                            @change="handleTableChange"
                            @expand="expandedRowsChange"
                    >
                        <a-table
                                slot="expandedRowRender"
                                slot-scope="text"
                                :columns="innerColumns"
                                :data-source="innerLog[text.id.toString()]"
                                :rowKey="r => r.id.toString()"
                                :pagination="false"
                        >
                            <span slot="id" slot-scope="text">{{text.toString()}}</span>
                            <span slot="time" slot-scope="text">{{DateTimeFormat(text)}}</span>

                        </a-table>

                        <span slot="id" slot-scope="text">{{text.toString()}}</span>
                        <span slot="time" slot-scope="text">{{DateTimeFormat(text)}}</span>

                    </a-table>
                </a-descriptions-item>
            </a-descriptions>
        </div>
        <div>
            <a-button type="primary" @click="back()">返回</a-button>
        </div>
    </div>
</template>

<script>

    import axios from "../../library/http";
    import {momentDate, DateTimeFormat} from "../../library/utils";

    export default {
        name: "FlowDetail",

        mounted() {
            this.fetch();
            this.fetchLogFlow();
        },
        data() {
            return {
                data: {},
                dataLogFlow: [],
                columns: [
                    {
                        title: 'Id',
                        dataIndex: 'id',
                        width: 200,
                        scopedSlots: {
                            customRender: 'id'
                        },
                    },
                    {
                        title: '工作流',
                        dataIndex: 'flow_id',
                        width: 200,
                    },
                    {
                        title: '开始时间',
                        dataIndex: 'start_time',
                        scopedSlots: {
                            customRender: 'time'
                        },
                    },
                    {
                        title: '结束时间',
                        dataIndex: 'end_time',
                        scopedSlots: {
                            customRender: 'time'
                        },
                    },
                ],
                innerColumns: [
                    {
                        title: 'Id',
                        dataIndex: 'id',
                        width: 180,
                        scopedSlots: {
                            customRender: 'id'
                        },
                    },
                    {
                        title: '任务',
                        dataIndex: 'job_id',
                        width: 200,
                    },

                    {
                        title: '输入',
                        dataIndex: 'input',
                        width: 200,
                    },

                    {
                        title: '输出',
                        dataIndex: 'output',
                        width: 200,
                    },

                    {
                        title: '错误信息',
                        dataIndex: 'error_msg',
                        width: 200,
                    },
                    {
                        title: '开始时间',
                        dataIndex: 'start_time',
                        scopedSlots: {
                            customRender: 'time'
                        },
                    },
                    {
                        title: '结束时间',
                        dataIndex: 'end_time',
                        scopedSlots: {
                            customRender: 'time'
                        },
                    },
                ],
                pagination: {},
                loading: false,
                innerLog: Map,
            };
        },
        methods: {
            fetch() {
                axios({
                    url: '/flow/detail',
                    data: {
                        id: this.$route.query.id,
                    },
                    type: 'json',
                }).then(data => {
                    console.log(data)

                    this.data = data.data;

                });
            },
            back() {
                if (this.$route.query.goindex === 'true') {
                    this.$router.push('/')
                } else {
                    this.$router.go(-1)
                }
            },
            momentDate,
            DateTimeFormat,
            handleTableChange(pagination) {
                console.log(pagination);
                const pager = {...this.pagination};
                pager.current = pagination.current;
                this.pagination = pager;
                this.fetchLogFlow({
                    page: pagination.current,
                });
            },
            fetchLogFlow(params = {}) {
                console.log('params:', params);
                this.loading = true;
                let offset = 0;
                if (params.page > 0) {
                    offset = (params.page - 1) * 10
                }
                axios({
                    url: '/log/flow/list',
                    data: {
                        flow_id: this.$route.query.id,
                        sort_value: "id desc",
                        offset: offset,
                        limit: 10,
                        ...params,
                    },
                    type: 'json',
                }).then(data => {
                    console.log(data)
                    const pagination = {...this.pagination};
                    // Read total count from server
                    // pagination.total = data.totalCount;
                    pagination.total = data.data.total;
                    this.loading = false;
                    data.data.data.forEach((value, index) => {
                        this.innerLog[value.id.toString()] = data.data.data[index].log_job;
                    })
                    this.dataLogFlow = data.data.data;
                    this.pagination = pagination;
                });
            },
            expandedRowsChange(expanded, record) {
                console.log(expanded, record)
            }
        },

    }
</script>

<style scoped>
    .highlight {
        background-color: rgb(255, 192, 105);
        padding: 0px;
    }
</style>