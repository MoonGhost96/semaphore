<template>
  <div class="auth">
    <v-dialog v-model="loginHelpDialog" max-width="600">
      <v-card>
        <v-card-title>
          How to fix sign-in issues
          <v-spacer></v-spacer>
          <v-btn icon @click="loginHelpDialog = false">
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>
        <v-card-text>
          <p class="text-body-1">
            Firstly, you need access to the server where Semaphore running.
          </p>
          <p class="text-body-1">
            Execute the following command on the server to see existing users:
          </p>
          <v-alert
            dense
            text
            color="info"
            style="font-family: monospace;"
          >
            semaphore user list
          </v-alert>
          <p class="text-body-1">
            You can change password of existing user:
          </p>
          <v-alert
            dense
            text
            color="info"
            style="font-family: monospace;"
          >
            semaphore user change-by-login --login user123 --password {{ makePasswordExample() }}
          </v-alert>
          <p class="text-body-1">
            Or create new admin user:
          </p>
          <v-alert
            dense
            text
            color="info"
            style="font-family: monospace;"
          >
            semaphore user add --admin --login user123 --name User123
            --email user123@example.com --password {{ makePasswordExample() }}
          </v-alert>
        </v-card-text>
        <v-card-actions>
          <v-spacer/>
          <v-btn
            color="blue darken-1"
            text
            @click="loginHelpDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-container
      fluid
      fill-height
      align-center
      justify-center
      class="pa-0"
    >
      <v-form
        ref="signInForm"
        lazy-validation
        v-model="signInFormValid"
        style="width: 300px; height: 300px;"
      >
        <h2 class="text-center mb-8">ANSIBLE-SEMAPHORE</h2>
        <h2 class="text-center mb-8">自动化运维工具</h2>

        <v-alert
          :value="signInError != null"
          color="error"
          style="margin-bottom: 20px;"
        >{{ signInError }}
        </v-alert>

        <v-text-field
          readonly
          onfocus="this.removeAttribute('readonly')"
          v-model="username"
          label="用户名"
          :rules="[v => !!v || '用户名为必填项']"
          required
          :disabled="signInProcess"
        ></v-text-field>

        <v-text-field
          v-model="password"
          label="密码"
          :rules="[v => !!v || '密码为必填项']"
          type="password"
          required
          :disabled="signInProcess"
          @keyup.enter.native="signIn"
          style="margin-bottom: 20px;"
        ></v-text-field>

        <v-btn
          color="primary"
          @click="signIn"
          :disabled="signInProcess"
          block
        >
          登录
        </v-btn>

        <div class="text-center mt-6">
          <a @click="loginHelpDialog = true" style="color: cornflowerblue">无法登录?</a>
        </div>
      </v-form>
    </v-container>
  </div>
</template>
<style lang="scss">
.auth {
  height: 100vh;
  background-image: url("~@/assets/background.jpg");
  background-size: cover;
  background-repeat: no-repeat;
}
</style>
<script>
import axios from 'axios';
import { getErrorMessage } from '@/lib/error';

export default {
  data() {
    return {
      signInFormValid: false,
      signInError: null,
      signInProcess: false,

      password: null,
      username: null,

      loginHelpDialog: null,
    };
  },

  async created() {
    if (this.isAuthenticated()) {
      document.location = document.baseURI;
    }
  },

  methods: {
    makePasswordExample() {
      let pwd = '';
      const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
      const charactersLength = characters.length;
      for (let i = 0; i < 10; i += 1) {
        pwd += characters.charAt(Math.floor(Math.random() * charactersLength));
      }
      return pwd;
    },

    isAuthenticated() {
      return document.cookie.includes('semaphore=');
    },

    async signIn() {
      this.signInError = null;

      if (!this.$refs.signInForm.validate()) {
        return;
      }

      this.signInProcess = true;
      try {
        await axios({
          method: 'post',
          url: '/api/auth/login',
          responseType: 'json',
          data: {
            auth: this.username,
            password: this.password,
          },
        });
        document.location = document.baseURI;
      } catch (err) {
        if (err.response.status === 401) {
          this.signInError = 'Incorrect login or password';
        } else {
          this.signInError = getErrorMessage(err);
        }
      } finally {
        this.signInProcess = false;
      }
    },
  },
};
</script>
