<template>

  <div>
    <div class="flex1 ">
      <div>

        <!-- <div v-for="(key, index) in data" :key="index">
          {{ key.name }}
        </div> -->

      </div>

      <div class="tasklist ">
        <div>
          <el-button @click="onMoreTask()">新增任务</el-button>

          <!-- <el-button @click="onRuleCreate"><i class="el-icon-setting "></i></el-button> -->

        </div>

        <div>
          <!-- <el-tag type="success">{{ $t('business.cluster.ready') }}</el-tag> -->
          <div v-for="(key, index) in worklist" :key="index">

            <div>
              <input type="checkbox" @click="selectTask(index)" checked="key.selected">aa</input>

              {{ key.name }}
            </div>


          </div>

        </div>

      </div>

    </div>

  </div>

</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import { deleteCluster, listClusters, searchClusters, updateCluster } from "@/api/clusters"
import { getTaskAll, setTask, getConfig, setConfig } from "@/api/capture"

import { checkPermissions } from "@/utils/permission";
import ComplexTable from "@/components/complex-table";
import Rule from "@/utils/rules"
import { downloadHelmReleases } from "@/utils/helm"

export default {
  name: "ClusterList",
  components: { LayoutContent, ComplexTable },
  data() {
    return {
      loading: false,
      guideDialogVisible: false,
      showAddLabelVisible: false,
      items: [],
      timer: null,
      // worklist: [], // 任务列表
      worklist: [
        {
          name: "taskaaa",
          selected: "no"
        },
        {
          name: "taskbbb",
          selected: "no"
        }
      ], // 任务列表
      currentConfig: null,
      // currentConfig: {
      //   config: {
      //     default_img: "bbb:v1",
      //     arr_img: [
      //       {
      //         cluster_name: "aaa",
      //         img: "bg4:v1",
      //       },
      //       {
      //         cluster_name: "ddd",
      //         img: "sh1:v1",
      //       },
      //     ]
      //   },
      //   uuid: ""
      // },
      data: [], // 集群信息
      datapods: [], // 集群所有pod信息-node信息
      selects: [],
      labelRules: {
        key: [
          Rule.RequiredRule,
          {
            min: 0,
            max: 10,
            message: this.$t("commons.validate.limit", [1, 10]),
            trigger: "blur"
          }
        ]
      },
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      }
    }
  },
  methods: {
    onMoreTask() {
      console.log("onMoreTask")
    },
    selectTask(index) {
      console.log(index)
    },
    search(conditions) {
      this.loading = true
      const { currentPage, pageSize } = this.paginationConfig
      searchClusters(currentPage, pageSize, conditions).then(data => {
        this.loading = false
        this.data = data.data.items
        console.log(this.data)
        this.paginationConfig.total = data.data.total
        this.data.forEach(d => {
          this.$set(d, "showAddLabelVisible", false)
          this.$set(d, "k", 0)
          this.$set(d, "form", {
            key: "",
          },)
        })
      })
    },
    pullingClusterStatus() {
      this.timer = setInterval(() => {
        listClusters().then(data => {
          this.items = data.data;
        })
      }, 3000)
    },
    //导出选中集群的所有原始helm releases
    async onExportAllHelmReleases() {
      this.loading = true
      await downloadHelmReleases(this.selects)
      this.loading = false
    }
  },
  created() {
    /*
    setTask({
      time_during: 88,
      items: [
        {
          podip: "192.168.113.12",
          nodeip: "192.168.113.12",
        }
      ],
      commands: [
        {
          command_line: "ls;sleep 9999",
          to_file: ""
        }
      ]
    }).then(data => {
      console.log(data.data)
      console.log(JSON.stringify(data.data))
    })
      */


    // setTask({}).then(data => {
    //   console.log(data.data)
    //   console.log(JSON.stringify(data.data))
    // })

    getTaskAll().then(data => {
      console.log(data.data)
      console.log(JSON.stringify(data.data))
    })

    // getConfig()

    /*
    getConfig().then(data => {
      console.log(data.data)
      console.log(JSON.stringify(data.data))
      this.currentConfig = data.data
      // this.currentConfig.config.default_img = "aacc:v5"
      // setConfig(this.currentConfig)
  })
  */
    /*
  let a = {
    config: {
      default_img: "bbb:v1",
      arr_img: [
        {
          cluster_name: "aaa",
          img: "bg4:v1",
        },
        {
          cluster_name: "ddd",
          img: "sh1:v1",
        },
      ]
    }
  }
  console.log(a)
  console.log(a)
  console.log(JSON.stringify(a))
  setConfig(a)
  */

  },
  destroyed() {
    clearInterval(this.timer)
  }
}
</script>

<style scoped>
.flex1 {
  /*flex布局（作用于容器）*/
  display: flex;

  /*两端对齐（作用于容器）*/
  justify-content: space-between;

  min-height: 400px;
  background-color: white;
}

.tasklist {
  background-color: rgb(226, 226, 226);

}

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}

.clearfix:after {
  clear: both
}

.bottom {
  margin-top: 13px;
  line-height: 12px;
}

.bottom-button {
  padding: 0;
  float: right;
}

.cluster-card {
  margin-left: 10px;
  margin-top: 20px;
  min-height: 169px;
}
</style>