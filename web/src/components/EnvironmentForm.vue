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
    >{{ formError }}</v-alert>

    <v-text-field
      v-model="item.name"
      label="环境名称"
      :rules="[v => !!v || '名称为必填项']"
      required
      :disabled="formSaving"
      class="mb-4"
    ></v-text-field>

    <v-subheader>
      Extra variables
    </v-subheader>

    <codemirror
        :style="{ border: '1px solid lightgray' }"
        v-model="item.json"
        :options="cmOptions"
        placeholder="Enter extra variables JSON..."
    />

    <v-subheader>
      Environment variables
    </v-subheader>

    <codemirror
        :style="{ border: '1px solid lightgray' }"
        v-model="item.env"
        :options="cmOptions"
        placeholder="Enter env JSON..."
    />

    <v-alert
        dense
        text
        type="info"
        class="mt-4"
    >
      Environment and extra variables must be valid JSON.
      Example:
      <pre style="font-size: 14px;">{
  "var_available_in_playbook_1": 1245,
  "var_available_in_playbook_2": "test"
}</pre>
    </v-alert>
  </v-form>
</template>
<script>
/* eslint-disable import/no-extraneous-dependencies,import/extensions */

import ItemFormBase from '@/components/ItemFormBase';

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
        mode: 'application/json',
        lineNumbers: true,
        line: true,
        lint: true,
        indentWithTabs: false,
      },
    };
  },

  methods: {
    getItemsUrl() {
      return `/api/project/${this.projectId}/environment`;
    },

    getSingleItemUrl() {
      return `/api/project/${this.projectId}/environment/${this.itemId}`;
    },
  },
};
</script>
