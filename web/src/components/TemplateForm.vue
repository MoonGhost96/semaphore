<template>
  <v-form
    ref="form"
    lazy-validation
    v-model="formValid"
    v-if="isLoaded"
  >
    <v-dialog
      v-model="helpDialog"
      hide-overlay
      width="300"
    >
      <v-alert
        border="top"
        colored-border
        type="info"
        elevation="2"
        class="mb-0 pb-0"
      >
        <div v-if="helpKey === 'build_version'">
          <p>
            Defines start version of your artifact.
            Each run increments the artifact version.
          </p>
          <p>
            For more information about building, see the
            <a href="https://docs.ansible-semaphore.com/user-guide/task-templates#build"
               target="_blank"
            >Task Template reference</a>.
          </p>
        </div>
        <div v-else-if="helpKey === 'build'">
          <p>
            Defines what artifact should be deployed when the task run.
          </p>
          <p>
            For more information about deploying, see the
            <a href="https://docs.ansible-semaphore.com/user-guide/task-templates#build"
               target="_blank"
            >Task Template reference</a>.
          </p>
        </div>
        <div v-if="helpKey === 'cron'">
          <p>Defines autorun schedule.</p>
          <p>
            For more information about cron, see the
            <a href="https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format"
               target="_blank"
            >Cron expression format reference</a>.
          </p>
        </div>
      </v-alert>
    </v-dialog>

    <v-alert
      :value="formError"
      color="error"
      class="pb-2"
    >{{ formError }}
    </v-alert>

    <v-row>
      <v-col cols="12" md="6" class="pb-0">
        <v-card class="mb-6" :color="$vuetify.theme.dark ? '#212121' : 'white'">
          <v-tabs
            fixed-tabs
            v-model="itemTypeIndex"
          >
            <v-tab
              style="padding: 0"
              v-for="(key) in Object.keys(TEMPLATE_TYPE_ICONS)"
              :key="key"
            >
              <v-icon small class="mr-2">{{ TEMPLATE_TYPE_ICONS[key] }}</v-icon>
              {{ TEMPLATE_TYPE_TITLES[key] }}
            </v-tab>
          </v-tabs>

          <div class="ml-4 mr-4 mt-6" v-if="item.type">
            <v-text-field
              v-if="item.type === 'build'"
              v-model="item.start_version"
              label="Start Version"
              :rules="[v => !!v || 'Start Version 为必填项']"
              required
              :disabled="formSaving"
              placeholder="Example: 0.0.0"
              append-outer-icon="mdi-help-circle"
              @click:append-outer="showHelpDialog('build_version')"
            ></v-text-field>

            <v-autocomplete
              v-if="item.type === 'deploy'"
              v-model="item.build_template_id"
              label="Build Template"
              :items="buildTemplates"
              item-value="id"
              item-text="name"
              :rules="[v => !!v || 'Build Template 为必填项']"
              required
              :disabled="formSaving"
              append-outer-icon="mdi-help-circle"
              @click:append-outer="showHelpDialog('build')"
            ></v-autocomplete>

            <v-checkbox
              v-if="item.type === 'deploy'"
              class="mt-0"
              label="Autorun"
              v-model="item.autorun"
            />
          </div>

        </v-card>

        <v-text-field
          v-model="item.name"
          label="名称*"
          :rules="[v => !!v || '名称为必填项']"
          outlined
          dense
          required
          :disabled="formSaving"
        ></v-text-field>

        <v-textarea
          v-model="item.description"
          label="描述"
          :disabled="formSaving"
          rows="1"
          :auto-grow="true"
          outlined
          dense
        ></v-textarea>

        <v-text-field
          v-if="item.type !== 'command'"
          v-model="item.playbook"
          label="Playbook文件路径*"
          :rules="[v => !!v || 'Playbook文件路径为必填项']"
          outlined
          dense
          required
          :disabled="formSaving"
          placeholder="Example: site.yml"
        ></v-text-field>

        <v-select
          v-model="item.inventory_id"
          label="资产目录*"
          :items="inventory"
          item-value="id"
          item-text="name"
          :rules="[v => !!v || '资产目录为必填项']"
          outlined
          dense
          required
          :disabled="formSaving"
        ></v-select>

        <v-select
          v-if="item.type !== '' && item.type !== 'command'"
          v-model="item.repository_id"
          label="仓库*"
          :items="repositories"
          item-value="id"
          item-text="name"
          :rules="[v => !!v || '仓库为必填项']"
          outlined
          dense
          required
          :disabled="formSaving"
        ></v-select>

        <v-select
          v-if="item.type === '' || item.type === 'command'"
          v-model="item.repository_id"
          label="仓库"
          :items="repositories"
          item-value="id"
          item-text="name"
          outlined
          dense
          :disabled="formSaving"
          clearable
        ></v-select>

        <v-select
          v-if="item.type !== '' && item.type !== 'command'"
          v-model="item.environment_id"
          label="环境变量*"
          :items="environment"
          item-value="id"
          item-text="name"
          :rules="[v => !!v || '环境变量为必填项']"
          outlined
          dense
          required
          :disabled="formSaving"
        ></v-select>

        <v-select
          v-if="item.type === ''"
          v-model="item.environment_id"
          label="环境变量"
          :items="environment"
          item-value="id"
          item-text="name"
          outlined
          dense
          :disabled="formSaving"
          clearable
        ></v-select>

<!--        <v-select-->
<!--          v-if="item.type === ''"-->
<!--          v-model="item.vault_key_id"-->
<!--          label="Ansible-Vault密码"-->
<!--          clearable-->
<!--          :items="loginPasswordKeys"-->
<!--          item-value="id"-->
<!--          item-text="name"-->
<!--          :disabled="formSaving"-->
<!--          outlined-->
<!--          dense-->
<!--        ></v-select>-->

        <v-select
          v-if="item.type === 'command'"
          v-model="item.module"
          label="命令模块*"
          :items="TEMPLATE_COMMAND_MODULES"
          :rules="[v => !!v || '命令模块为必填项']"
          outlined
          dense
          required
          :disabled="formSaving"
        ></v-select>

        <codemirror
          v-if="item.type === 'command'"
          :style="{ border: '1px solid lightgray' }"
          v-model="item.command"
          :options="commandCmOptions"
          :disabled="formSaving"
          placeholder="命令:
ls -l >> output.txt
echo 'hello world'"
        />

      </v-col>

      <v-col cols="12" md="6" class="pb-0">

<!--        <v-select-->
<!--          v-if="item.type !== ''"-->
<!--          v-model="item.vault_key_id"-->
<!--          label="Ansible-Vault密码"-->
<!--          clearable-->
<!--          :items="loginPasswordKeys"-->
<!--          item-value="id"-->
<!--          item-text="name"-->
<!--          :disabled="formSaving"-->
<!--          outlined-->
<!--          dense-->
<!--        ></v-select>-->

        <SurveyVars style="margin-top: -10px;" :vars="item.survey_vars" @change="setSurveyVars"/>

        <v-select
          v-model="item.view_id"
          label="视图"
          clearable
          :items="views"
          item-value="id"
          item-text="title"
          :disabled="formSaving"
          outlined
          dense
        ></v-select>

        <v-row>
          <v-col cols="5" class="pr-1">
            <v-text-field
              style="font-size: 14px"
              v-model="cronFormat"
              label="Cron定时任务"
              :disabled="formSaving"
              placeholder="* * * * *"
              v-if="schedules == null || schedules.length <= 1"
              outlined
              dense
              hide-details
            ></v-text-field>
          </v-col>

          <v-col
            cols="7"
            v-if="item.type !== 'command'"
          >
            <a
              v-if="!cronRepositoryIdVisible && cronRepositoryId == null"
              @click="cronRepositoryIdVisible = true"
              class="text-caption d-block"
              style="line-height: 1.1;"
            >
              针对某个版本仓库的新commit通过cron定时执行一个任务。
            </a>

            <v-select
              style="font-size: 14px"
              v-if="cronRepositoryIdVisible || cronRepositoryId != null"
              v-model="cronRepositoryId"
              label="Repository"
              placeholder="Cron checks new commit before run"
              :items="repositories"
              item-value="id"
              item-text="name"
              clearable
              :disabled="formSaving"
              outlined
              dense
              hide-details
            ></v-select>

          </v-col>
        </v-row>

        <small class="mt-1 mb-4 d-block">
          请使用五位Cron格式，例如"* * * * *"。更多信息请参考
          <a target="_blank" href="https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format">Cron文档</a>
        </small>

        <v-checkbox
          class="mt-0"
          label="运行成功无需告警"
          v-model="item.suppress_success_alerts"
        />

<!--        <a @click="advancedOptions = true" v-if="!advancedOptions">-->
<!--          Advanced-->
<!--          <v-icon style="transform: translateY(-1px)">mdi-chevron-right</v-icon>-->
<!--        </a>-->

<!--        <div v-if="advancedOptions" class="mb-3">-->
<!--          <a @click="advancedOptions = false">-->
<!--            Hide-->
<!--            <v-icon style="transform: translateY(-1px)">mdi-chevron-up</v-icon>-->
<!--          </a>-->
<!--        </div>-->

        <codemirror
          :style="{ border: '1px solid lightgray' }"
          v-model="item.arguments"
          :options="jsonCmOptions"
          :disabled="formSaving"
          placeholder='自定义运行参数(JSON格式)，例如:
[
  "-i",
  "@myinventory.sh",
  "--private-key=/there/id_rsa",
  "-vvvv"
]'
        />

        <v-checkbox
          label="允许在执行时指定参数"
          v-model="item.allow_override_args_in_task"
        />

      </v-col>
    </v-row>
  </v-form>
</template>
<style lang="scss">
.CodeMirror-placeholder {
  color: #a4a4a4 !important;
}
</style>
<script>
/* eslint-disable import/no-extraneous-dependencies,import/extensions */

import axios from 'axios';
// 引入CodeMirror和基础样式
import { codemirror } from 'vue-codemirror';
import ItemFormBase from '@/components/ItemFormBase';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/vue/vue.js';
// 支持JSON格式
import 'codemirror/addon/lint/json-lint.js';
import 'codemirror/addon/display/placeholder.js';
// 明亮和暗黑主题
import 'codemirror/theme/base16-dark.css';
import 'codemirror/theme/base16-light.css';
// shell支持
import 'codemirror/mode/shell/shell.js';
// JSON代码高亮需要由JavaScript插件支持
import 'codemirror/mode/javascript/javascript.js';
// 支持括号自动匹配
import 'codemirror/addon/edit/matchbrackets.js';
import 'codemirror/addon/edit/closebrackets.js';
// 行注释
import 'codemirror/addon/comment/comment.js';
// JSON错误检查
import 'codemirror/addon/lint/lint.css';
import 'codemirror/addon/lint/lint.js';
// 搜索功能的依赖
import 'codemirror/addon/dialog/dialog.js';
import 'codemirror/addon/dialog/dialog.css';
// 支持搜索功能
import 'codemirror/addon/search/search';
import 'codemirror/addon/search/searchcursor.js';
// 支持代码区域全屏功能
import 'codemirror/addon/display/fullscreen.css';
import 'codemirror/addon/display/fullscreen.js';
import { TEMPLATE_TYPE_ICONS, TEMPLATE_TYPE_TITLES, TEMPLATE_COMMAND_MODULES } from '@/lib/constants';
import SurveyVars from './SurveyVars';

export default {
  mixins: [ItemFormBase],

  components: {
    SurveyVars,
    codemirror,
  },

  props: {
    sourceItemId: Number,
  },

  data() {
    return {
      itemTypeIndex: 0,
      TEMPLATE_TYPE_ICONS,
      TEMPLATE_TYPE_TITLES,
      TEMPLATE_COMMAND_MODULES,
      item: null,
      keys: null,
      inventory: null,
      repositories: null,
      environment: null,
      views: null,
      schedules: null,
      buildTemplates: null,
      cronFormat: null,
      cronRepositoryId: null,
      cronRepositoryIdVisible: false,

      helpDialog: null,
      helpKey: null,

      advancedOptions: false,
    };
  },

  watch: {
    needReset(val) {
      if (val) {
        if (this.item != null) {
          this.item.template_id = this.templateId;
        }
      }
    },

    sourceItemId(val) {
      this.item.template_id = val;
    },

    itemTypeIndex(val) {
      this.item.type = Object.keys(TEMPLATE_TYPE_ICONS)[val];
    },
  },

  computed: {
    jsonCmOptions() {
      return {
        tabSize: 2,
        mode: 'application/json',
        styleActiveLine: true,
        lineNumbers: true,
        line: true,
        theme: this.$vuetify.theme.dark ? 'base16-dark' : 'base16-light',
        // 括号匹配
        matchBrackets: true,
        autoCloseBrackets: true,
      };
    },

    commandCmOptions() {
      return {
        smartIndent: true,
        tabSize: 2,
        mode: 'shell',
        styleActiveLine: true,
        lineNumbers: true,
        line: true,
        lint: true,
        indentWithTabs: true,
        theme: this.$vuetify.theme.dark ? 'base16-dark' : 'base16-light',
        // 括号匹配
        matchBrackets: true,
        autoCloseBrackets: true,
      };
    },

    isLoaded() {
      if (this.isNew && this.sourceItemId == null) {
        return true;
      }

      return this.keys != null
        && this.repositories != null
        && this.inventory != null
        && this.environment != null
        && this.item != null
        && this.schedules != null
        && this.views != null;
    },

    loginPasswordKeys() {
      if (this.keys == null) {
        return null;
      }
      return this.keys.filter((key) => key.type === 'login_password');
    },
  },

  methods: {
    setSurveyVars(v) {
      this.item.survey_vars = v;
    },

    showHelpDialog(key) {
      this.helpKey = key;
      this.helpDialog = true;
    },

    async afterLoadData() {
      if (this.sourceItemId) {
        this.item = (await axios({
          keys: 'get',
          url: `/api/project/${this.projectId}/templates/${this.sourceItemId}`,
          responseType: 'json',
        })).data;
      }

      this.advancedOptions = this.item.arguments != null || this.item.allow_override_args_in_task;

      this.keys = (await axios({
        keys: 'get',
        url: `/api/project/${this.projectId}/keys`,
        responseType: 'json',
      })).data;

      this.repositories = (await axios({
        keys: 'get',
        url: `/api/project/${this.projectId}/repositories`,
        responseType: 'json',
      })).data;

      this.inventory = (await axios({
        keys: 'get',
        url: `/api/project/${this.projectId}/inventory`,
        responseType: 'json',
      })).data;

      this.environment = (await axios({
        keys: 'get',
        url: `/api/project/${this.projectId}/environment`,
        responseType: 'json',
      })).data;

      const template = (await axios({
        keys: 'get',
        url: `/api/project/${this.projectId}/templates`,
        responseType: 'json',
      })).data;
      const builds = [];
      const deploys = [];
      template.forEach((t) => {
        switch (t.type) {
          case 'build':
            if (builds.length === 0) {
              builds.push({ header: 'Build Templates' });
            }
            builds.push(t);
            break;
          case 'deploy':
            if (deploys.length === 0) {
              deploys.push({ header: 'Deploy Templates' });
            }
            deploys.push(t);
            break;
          default:
            break;
        }
      });

      this.buildTemplates = builds;
      if (this.buildTemplates.length > 0 && deploys.length > 0) {
        this.buildTemplates.push({ divider: true });
      }
      this.buildTemplates.push(...deploys);

      this.schedules = this.isNew ? [] : (await axios({
        keys: 'get',
        url: `/api/project/${this.projectId}/templates/${this.itemId}/schedules`,
        responseType: 'json',
      })).data;

      this.views = (await axios({
        keys: 'get',
        url: `/api/project/${this.projectId}/views`,
        responseType: 'json',
      })).data;

      if (this.schedules.length === 1) {
        this.cronFormat = this.schedules[0].cron_format;
        this.cronRepositoryId = this.schedules[0].repository_id;
      }

      this.itemTypeIndex = Object.keys(TEMPLATE_TYPE_ICONS).indexOf(this.item.type);
    },

    getItemsUrl() {
      return `/api/project/${this.projectId}/templates`;
    },

    getSingleItemUrl() {
      return `/api/project/${this.projectId}/templates/${this.itemId}`;
    },

    async beforeSave() {
      if (this.cronFormat == null || this.cronFormat === '') {
        return;
      }

      await axios({
        method: 'post',
        url: `/api/project/${this.projectId}/schedules/validate`,
        responseType: 'json',
        data: {
          cron_format: this.cronFormat,
        },
      });
    },

    async afterSave(newItem) {
      if (newItem || this.schedules.length === 0) {
        if (this.cronFormat != null && this.cronFormat !== '') {
          // new schedule
          await axios({
            method: 'post',
            url: `/api/project/${this.projectId}/schedules`,
            responseType: 'json',
            data: {
              project_id: this.projectId,
              template_id: newItem ? newItem.id : this.itemId,
              cron_format: this.cronFormat,
              repository_id: this.cronRepositoryId,
            },
          });
        }
      } else if (this.schedules.length > 1) {
        // do nothing
      } else if (this.cronFormat == null || this.cronFormat === '') {
        // drop schedule
        await axios({
          method: 'delete',
          url: `/api/project/${this.projectId}/schedules/${this.schedules[0].id}`,
          responseType: 'json',
        });
      } else {
        // update schedule
        await axios({
          method: 'put',
          url: `/api/project/${this.projectId}/schedules/${this.schedules[0].id}`,
          responseType: 'json',
          data: {
            id: this.schedules[0].id,
            project_id: this.projectId,
            template_id: this.itemId,
            cron_format: this.cronFormat,
            repository_id: this.cronRepositoryId,
          },
        });
      }
    },
  },
};
</script>
