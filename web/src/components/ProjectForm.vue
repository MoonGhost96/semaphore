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
      label="Project Name"
      :rules="[v => !!v || 'Project name 为必填项']"
      required
      :disabled="formSaving"
    ></v-text-field>

    <v-checkbox
      v-model="item.alert"
      label="允许该项目发出告警"
    ></v-checkbox>

    <v-text-field
      v-model="item.alert_chat"
      label="Telegram Chat ID (Optional)"
      :disabled="formSaving"
    ></v-text-field>

    <v-text-field
      v-model.number="item.max_parallel_tasks"
      label="Max number of parallel tasks (Optional)"
      :disabled="formSaving"
      :rules="[
        v => (v == null || v === '' || Math.floor(v) === v) || 'Must be integer',
        v => (v == null || v === '' || v >= 0) || 'Must be 0 or greater',
      ]"
      hint="Should be 0 or greater, 0 - unlimited."
      type="number"
      :step="1"
    ></v-text-field>
  </v-form>
</template>
<script>
import ItemFormBase from '@/components/ItemFormBase';

export default {
  mixins: [ItemFormBase],
  methods: {
    getItemsUrl() {
      return '/api/projects';
    },
    getSingleItemUrl() {
      return `/api/project/${this.itemId}`;
    },
  },
};
</script>
