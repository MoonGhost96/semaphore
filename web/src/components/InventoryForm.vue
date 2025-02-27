<template>
  <v-form
    ref="form"
    lazy-validation
    v-model="formValid"
    v-if="item != null && keys != null"
  >
    <v-alert
      :value="formError"
      color="error"
      class="pb-2"
    >{{ formError }}</v-alert>

    <v-text-field
      v-model="item.name"
      label="名称"
      :rules="[v => !!v || '名称为必填项']"
      required
      :disabled="formSaving"
    ></v-text-field>

    <v-select
      v-model="item.ssh_key_id"
      label="用户凭证"
      :items="keys"
      item-value="id"
      item-text="name"
      :rules="[v => !!v || '用户凭证为必填项']"
      required
      :disabled="formSaving"
    ></v-select>

    <v-select
        v-model="item.become_key_id"
        label="Sudo凭证(可选)"
        clearable
        :items="loginPasswordKeys"
        item-value="id"
        item-text="name"
        :disabled="formSaving"
    ></v-select>

    <v-select
      v-model="item.type"
      label="类型"
      :rules="[v => !!v || '类型为必填项']"
      :items="inventoryTypes"
      item-value="id"
      item-text="name"
      required
      :disabled="formSaving"
    ></v-select>

    <v-text-field
      v-model="item.inventory"
      label="Path to Inventory file"
      :rules="[v => !!v || 'Path to Inventory file 为必填项']"
      required
      :disabled="formSaving"
      v-if="item.type === 'file'"
    ></v-text-field>

    <v-checkbox
      v-model="item.host_inv_rels"
      label="Host List"
      v-if="item.type === 'host'"
    />

    <codemirror
        :style="{ border: '1px solid lightgray' }"
        v-model="item.inventory"
        :options="cmOptions"
        v-if="item.type === 'static' || item.type === 'static-yaml'"
        placeholder="Enter inventory..."
    />

    <v-alert
        dense
        text
        class="mt-4"
        type="info"
        v-if="item.type === 'static'"
    >
      Static inventory example:
      <pre style="font-size: 14px;">[website]
172.18.8.40
172.18.8.41</pre>
    </v-alert>

    <v-alert
        dense
        text
        class="mt-4"
        type="info"
        v-if="item.type === 'static-yaml'"
    >
      Static YAML inventory example:
      <pre style="font-size: 14px;">all:
  children:
    website:
      hosts:
        172.18.8.40:
        172.18.8.41:</pre>
    </v-alert>
  </v-form>
</template>
<style>
.CodeMirror {
  height: 160px !important;
}
</style>
<script>
/* eslint-disable import/no-extraneous-dependencies,import/extensions */

import ItemFormBase from '@/components/ItemFormBase';
import axios from 'axios';

import { codemirror } from 'vue-codemirror';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/vue/vue.js';
import 'codemirror/addon/display/placeholder.js';

export default {
  mixins: [ItemFormBase],

  components: {
    codemirror,
  },

  data() {
    return {
      cmOptions: {
        tabSize: 2,
        mode: 'text/x-ini',
        lineNumbers: true,
        line: true,
        lint: true,
        indentWithTabs: false,
      },
      keys: null,
      inventoryTypes: [{
        id: 'static',
        name: 'Static',
      }, {
        id: 'static-yaml',
        name: 'Static YAML',
      }, {
        id: 'file',
        name: 'File',
      }, {
        id: 'host',
        name: 'Host',
      }],
    };
  },

  computed: {
    loginPasswordKeys() {
      if (this.keys == null) {
        return null;
      }
      return this.keys.filter((key) => key.type === 'login_password');
    },
  },

  async created() {
    this.keys = (await axios({
      keys: 'get',
      url: `/api/project/${this.projectId}/keys`,
      responseType: 'json',
    })).data;
  },

  methods: {
    getItemsUrl() {
      return `/api/project/${this.projectId}/inventory`;
    },
    getSingleItemUrl() {
      return `/api/project/${this.projectId}/inventory/${this.itemId}`;
    },
  },
};
</script>
