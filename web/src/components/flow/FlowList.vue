<template>
    <div>
        <h3 align="left">工作流 > 列表</h3>
        <a-form-model layout="inline" align="left" :style="{width:'100%'}">
            <a-form-model-item>
                <a-input placeholder="Id" v-model="id"></a-input>
            </a-form-model-item>
            <a-form-model-item>
                <a-input placeholder="名称" v-model="name"></a-input>
            </a-form-model-item>
            <a-form-model-item>
                <a-button type="primary" @click="filter">筛选</a-button>
            </a-form-model-item>
        </a-form-model>

        <a-table
                :columns="columns"
                :rowKey="record => record.id.toString()"
                :dataSource="data"
                :pagination="pagination"
                :loading="loading"
                @change="handleTableChange"
                :scroll="{ x:  'calc(860px + 40%)' }"
        >

            <span slot="time" slot-scope="text, record">{{DateTimeFormat(record.create_time)}}</span>

            <span slot="action" slot-scope="text, record">
                <a @click="() => detail(record.id.toString())">详情</a>
                <a-divider type="vertical"/>
                <a @click="() => doFlow(record.id.toString())">重试</a>
            </span>
        </a-table>

        <a-modal :closable="false"
                 centered
                 :visible="confirmModalShow"
                 @ok="doFlowNow"
                 @cancel="() => {this.confirmModalShow= false}"
                 cancelText="取消"
                 okText="确认"
                 title="重试">
            <a-form-item>
                <a-form-model :label-col="labelCol" :wrapper-col="wrapperCol">
                    <a-form-model-item label="工作流Id">
                        <label>{{flowId}}</label>
                    </a-form-model-item>
                    <a-form-model-item v-if="flowId === 'Environ'" label="日期">
                        <a-date-picker show-time v-model="flowDate"/>
                    </a-form-model-item>
                    <a-form-model-item label="重试任务">
                        <a-checkbox-group :options="options" @change="onCheckboxChange"/>
                    </a-form-model-item>
                </a-form-model>
            </a-form-item>
        </a-modal>

    </div>
</template>

<script>
    import axios from "../../library/http";
    import {DateTimeFormat} from "../../library/utils";

    export default {
        name: "FlowList",
        mounted() {
            this.fetch();
        },
        data() {
            return {
                dataClient: [],
                data: [],
                pagination: {},
                loading: false,
                columns: [
                    {
                        title: 'Id',
                        dataIndex: 'id',
                        width: 180,
                    },
                    {
                        title: '名称',
                        dataIndex: 'name',
                        width: 200,
                    },
                    {
                        title: '类型（1周期任务，2即时任务）',
                        dataIndex: 'flow_type',
                    },
                    {
                        title: '计划时间',
                        dataIndex: 'spec',
                        width: 180,
                    },
                    {
                        title: '1启用，2禁用',
                        dataIndex: 'status',
                        width: 120,
                    },
                    {
                        title: '创建时间',
                        dataIndex: 'create_time',
                        width: 180,
                        scopedSlots: {
                            customRender: 'time'
                        },
                    },
                    {
                        title: '操作',
                        key: 'action',
                        fixed: 'right',
                        width: 180,
                        scopedSlots: {
                            customRender: 'action'
                        },
                    },
                ],
                id: '',
                name: '',
                confirmModalShow: false,
                labelCol: {span: 4},
                wrapperCol: {span: 14},
                flowId: '',
                flowDate: '',
                flowJobs: [],
                options: [],
            };
        },
        methods: {
            handleTableChange(pagination) {
                console.log(pagination);
                const pager = {...this.pagination};
                pager.current = pagination.current;
                this.pagination = pager;
                this.fetch({
                    page: pagination.current,
                });
            },
            fetch(params = {}) {
                if (this.$route.query.id !== undefined && this.$route.query.id !== '') {
                    this.id = this.$route.query.id;
                }
                if (this.$route.query.name !== undefined && this.$route.query.name !== '') {
                    this.name = this.$route.query.name;
                }
                console.log('params:', params);
                this.loading = true;
                let offset = 0;
                if (params.page > 0) {
                    offset = (params.page - 1) * 10
                }
                axios({
                    url: '/flow/list',
                    data: {
                        id: this.id,
                        name: this.name,
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
                    this.data = data.data.data;
                    this.pagination = pagination;
                });
            },
            detail(id) {
                this.$router.push("/flow/detail?id=" + id);
            },
            filter() {
                let para = [];
                if (this.id !== undefined && this.id !== '') {
                    para = [...para, 'id=' + this.id]
                }
                if (this.name !== undefined && this.name !== '') {
                    para = [...para, 'name=' + this.name]
                }
                let url = "/flow/list";
                if (para.length > 0) {
                    url = url + "?" + para.join('&');
                }
                if (this.$route.fullPath !== url) {
                    this.$router.push(url);
                } else {
                    location.reload();
                }
            },
            DateTimeFormat,
            doFlow(flowId) {
                this.flowId = flowId;
                axios({
                    url: '/flow/detail',
                    data: {
                        id: flowId,
                    },
                    type: 'json',
                }).then(data => {
                    this.options = [];
                    const obj = JSON.parse(data.data.info);
                    const keys = Object.keys(obj.jobs);
                    for (let key in keys) {
                        this.options = [...this.options, {
                            label: "第" + keys[key] + "步",
                            value: keys[key]
                        }]
                    }
                    this.confirmModalShow = true;
                });
            },
            doFlowNow() {
                let m = {};
                for (let k in this.flowJobs) {
                    console.log(k)
                    console.log(this.flowJobs[k])
                    m[this.flowJobs[k]] = null;
                }

                console.log(m);

                axios({
                    url: '/flow/do',
                    data: {
                        flow_id: this.flowId,
                        date: Math.floor(this.flowDate.valueOf() / 1000),
                        do_force: m,
                    },
                    type: 'json',
                }).then(() => {
                    this.$message.info("任务流提交成功");
                    this.confirmModalShow = false;
                });
            },
            onCheckboxChange(checkedValues) {
                this.flowJobs = checkedValues;
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