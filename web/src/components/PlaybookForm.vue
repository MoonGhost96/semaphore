<template>
  <v-form
      ref="form"
      lazy-validation
      v-model="formValid"
      v-if="item != null"
  >
    <v-alert
        :value="formError"
        color="error"
        class="pb-2"
    >{{ formError }}
    </v-alert>

    <v-text-field
        v-model="item.name"
        label="剧本名称"
        :rules="[v => !!v || '剧本名称为必填项']"
        required
        :disabled="formSaving"
    />

    <v-text-field
        v-model="item.description"
        label="描述"
        :disabled="formSaving"
    />

    <codemirror
      class="mt-4"
      :style="{ border: '1px solid lightgray' }"
      v-model="item.content"
      :options="cmOptions"
      placeholder='Ansible playbook. Example:
---
- hosts: all (你应当将hosts设为all以使用资产目录来设定被控主机)
  tasks:
    - name: Install nginx
      apt:
        name: nginx
        state: present
    - name: Start nginx
      service:
        name: nginx
        state: started
'
    />

    <v-checkbox
        v-model="item.override_secret"
        label="Override"
        v-if="!isNew"
    />

    <v-alert
        dense
        text
        type="info"
        v-if="item.type === 'none'"
    >
      Use this type of key for HTTPS repositories and for
      playbooks which use non-SSH connections.
    </v-alert>
  </v-form>
</template>
<script>
/* eslint-disable import/no-extraneous-dependencies,import/extensions */
import { codemirror } from 'vue-codemirror';
import ItemFormBase from '@/components/ItemFormBase';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/vue/vue.js';
import 'codemirror/addon/lint/json-lint.js';
import 'codemirror/addon/lint/yaml-lint.js';
import 'codemirror/addon/display/placeholder.js';
import 'codemirror/theme/base16-dark.css';
import 'codemirror/theme/base16-light.css';
import 'codemirror/mode/yaml/yaml.js';

export default {
  components: {
    codemirror,
  },
  mixins: [ItemFormBase],
  data() {
    return {
      inventoryTypes: [{
        id: 'ssh',
        name: 'SSH Key',
      }, {
        id: 'login_password',
        name: 'Login with password',
      }, {
        id: 'none',
        name: 'None',
      }],
    };
  },

  computed: {
    canEditSecrets() {
      return this.isNew || this.item.override_secret;
    },

    cmOptions() {
      return {
        tabSize: 2,
        mode: 'yaml',
        lineNumbers: true,
        line: true,
        lint: true,
        indentWithTabs: false,
        theme: this.$vuetify.theme.dark ? 'base16-dark' : 'base16-light',
      };
    },
  },

  methods: {
    getNewItem() {
      return {
        ssh: {},
        login_password: {},
      };
    },

    getItemsUrl() {
      return `/api/project/${this.projectId}/playbooks`;
    },

    getSingleItemUrl() {
      return `/api/project/${this.projectId}/playbooks/${this.itemId}`;
    },
  },
};
</script>
