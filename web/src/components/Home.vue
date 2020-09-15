<template>
    <div>
        <a-layout id="components-layout-fixed">
            <a-layout-header :style="{ position: 'fixed', zIndex: 1, width: '100%' }">
                <div class="logo"/>
                <a-menu
                        theme="dark"
                        mode="horizontal"
                        :defaultSelectedKeys="['flow']"
                        :style="{ lineHeight: '64px' }"
                        :selectedKeys="selectedKeys"
                        @click="menuClick"
                >
                    <a-menu-item key="flow" class="menu-item">工作流</a-menu-item>
                    <a-menu-item key="job" class="menu-item">任务</a-menu-item>
                    <a-menu-item key="active" class="menu-item">活动</a-menu-item>
                </a-menu>
            </a-layout-header>
            <a-layout-content :style="{ padding: '0 25px', marginTop: '64px' }">
                <div :is="item.component" :style="{ background: '#fff', padding: '12px' }">
                </div>
            </a-layout-content>
            <a-layout-footer :style="{ textAlign: 'center' }">
                Copyright ©2020 Created by 杭州翰朗环保科技有限公司
            </a-layout-footer>
        </a-layout>
    </div>

</template>

<script>
    import System404 from "./404";

    import FlowDetail from "./flow/FlowDetail";
    import FlowList from "./flow/FlowList";
    import JobDetail from "./job/JobDetail";
    import JobList from "./job/JobList";
    import LogFlowDetail from "./log/LogFlowDetail";
    import LogFlowList from "./log/LogFlowList";
    import LogJobDetail from "./log/LogJobDetail";
    import LogJobList from "./log/LogJobList";
    import ActiveList from "./active/ActiveList";

    export default {
        name: "Home",
        data() {
            return {
                item: {},
                selectedKeys: [],
            };
        },
        mounted() {
            this.init();
        },
        components: {
            System404,
            FlowDetail,
            FlowList,
            JobList,
            JobDetail,
            LogJobList,
            LogJobDetail,
            LogFlowList,
            LogFlowDetail,
            ActiveList,
        },
        methods: {
            init() {
                const path = this.$route.path;
                console.log(this.$route.path);
                switch (path) {
                    case '/':
                    case '/flow/list' : {
                        this.item = {component: FlowList, text: '工作流'};
                        this.selectedKeys = ['flow'];
                        break;
                    }
                    case '/flow/detail' : {
                        this.item = {component: FlowDetail, text: '工作流'};
                        this.selectedKeys = ['flow'];
                        break;
                    }
                    case '/job/list' : {
                        this.item = {component: JobList, text: '任务'};
                        this.selectedKeys = ['job'];
                        break;
                    }
                    case '/job/detail' : {
                        this.item = {component: JobDetail, text: '任务'};
                        this.selectedKeys = ['job'];
                        break;
                    }
                    case '/log/flow/list' : {
                        this.item = {component: LogFlowList, text: '日志'};
                        this.selectedKeys = ['flow'];
                        break;
                    }
                    case '/log/flow/detail' : {
                        this.item = {component: LogFlowDetail, text: '日志'};
                        this.selectedKeys = ['flow'];
                        break;
                    }
                    case '/log/job/list' : {
                        this.item = {component: LogJobList, text: '日志'};
                        this.selectedKeys = ['job'];
                        break;
                    }
                    case '/log/job/detail' : {
                        this.item = {component: LogJobDetail, text: '日志'};
                        this.selectedKeys = ['job'];
                        break;
                    }
                    case '/active/list' : {
                        this.item = {component: ActiveList, text: '活动'};
                        this.selectedKeys = ['active'];
                        break;
                    }
                    default: {
                        this.item = {component: System404, text: '404'}
                    }
                }
            },
            menuClick(e) {
                const path = '/' + e.key + '/list';
                if (this.$route.path !== path) {
                    this.$router.push(path);
                }
            },
        },
    }
</script>

<style scoped>
    #components-layout-fixed .logo {
        width: 60px;
        height: 48px;
        background: rgba(255, 255, 255, 0.2);
        margin: 8px 12px 8px 0;
        float: left;
    }

    .menu-item {
        width: 120px;
        float: left;
    }
</style>