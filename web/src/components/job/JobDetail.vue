<template>
    <div>
        <h3 align="left">任务 > 详情</h3>
        <div align="left">
            <a-descriptions layout="vertical" :column="4" bordered>
                <a-descriptions-item label="Id">{{data.id}}</a-descriptions-item>
                <a-descriptions-item label="名称">{{data.name}}</a-descriptions-item>
                <a-descriptions-item label="类型">{{data.job_type}}</a-descriptions-item>
                <a-descriptions-item label="状态(1可用, 2禁用)">{{data.status}}</a-descriptions-item>
                <a-descriptions-item label="描述">{{data.desc}}</a-descriptions-item>
                <a-descriptions-item label="创建时间">{{DateTimeFormat(data.create_time)}}</a-descriptions-item>
                <a-descriptions-item label="更新时间">{{DateTimeFormat(data.update_time)}}</a-descriptions-item>
                <a-descriptions-item ></a-descriptions-item>
                <a-descriptions-item label="任务信息" :span="4">{{data.info}}</a-descriptions-item>
                <a-descriptions-item label="日志" :span="4">
                    <a-table
                            :columns="columns"
                            :data-source="dataLogJob"
                            :rowKey="record => record.id.toString()"
                            :pagination="pagination"
                            :loading="loading"
                            @change="handleTableChange"
                    >
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
        name: "JobDetail",

        mounted() {
            this.fetch();
            this.fetchLogJob();
        },
        data() {
            return {
                data: {},
                dataLogJob: [],
                columns: [
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
            };
        },
        methods: {
            fetch() {
                axios({
                    url: '/job/detail',
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
                this.fetchLogJob({
                    page: pagination.current,
                });
            },
            fetchLogJob(params = {}) {
                console.log('params:', params);
                this.loading = true;
                let offset = 0;
                if (params.page > 0) {
                    offset = (params.page - 1) * 10
                }
                axios({
                    url: '/log/job/list',
                    data: {
                        job_id: this.$route.query.id,
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
                    this.dataLogJob = data.data.data;
                    this.pagination = pagination;
                });
            },
        },

    }
</script>

<style scoped>
    .highlight {
        background-color: rgb(255, 192, 105);
        padding: 0px;
    }
</style>