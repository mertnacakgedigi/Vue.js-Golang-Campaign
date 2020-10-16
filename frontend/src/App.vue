<template>
  <v-app id="inspire">
    <v-data-table
      :headers="headers"
      :items="campaigns"
      sort-by="calories"
      class="elevation-1"
      :loading="isLoading"
      loading-text="Loading your campaigns ..."
    >
      <template v-slot:[`item.status`]="{ item }">
        <v-chip :color="getColor(item.status)" dark>
          {{ item.status }}
        </v-chip>
      </template>
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title>My Campaigns </v-toolbar-title>
          <v-divider class="mx-4" inset vertical></v-divider>
          <v-spacer></v-spacer>
          <v-dialog v-model="dialog" max-width="600px">
            <template v-slot:activator="{ on }">
              <v-btn color="primary" dark class="mb-2" v-on="on"
                >New Campaign</v-btn
              >
            </template>
            <v-card>
              <v-card-title>
                <span class="headline">{{ formTitle }}</span>
              </v-card-title>

              <v-card-text>
                <v-container>
                  <v-row>
                    <v-col cols="12" sm="6" md="4">
                      <v-text-field
                        v-model="editedItem.name"
                        label="Name"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" sm="6" md="4">
                      <v-select
                        v-model="editedItem.status"
                        :items="status"
                        label="Status"
                      ></v-select>
                    </v-col>
                    <v-col cols="12" sm="6" md="4">
                      <v-select
                        v-model="editedItem.type"
                        :items="type"
                        label="Type"
                      ></v-select>
                    </v-col>
                    <v-col cols="12" sm="6" md="4">
                      <v-text-field
                        v-model="editedItem.budget"
                        label="Budget"
                        prefix="$"
                      ></v-text-field>
                    </v-col>
                  </v-row>
                </v-container>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="blue darken-1" text @click="close">Cancel</v-btn>
                <v-btn color="blue darken-1" text @click="save">Save</v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
        </v-toolbar>
      </template>
      <template v-slot:[`item.action`]="{ item }">
        <v-icon small class="mr-2" @click="editItem(item.id)"> edit </v-icon>
        <v-icon small @click="deleteItem(item.id)"> delete </v-icon>
      </template>
      <template v-slot:no-data>
        <v-btn color="primary" @click="initialize">Reset</v-btn>
      </template>
    </v-data-table>
  </v-app>
</template>

<script>
import axios from "axios";
export default {
  data: () => ({
    dialog: false,
    headers: [
      {
        text: "Name",
        align: "left",
        sortable: false,
        value: "name",
      },
      { text: "Status", value: "status" },
      { text: "Type", value: "type" },
      { text: "Budget ($)", value: "budget" },
      { text: "Created", value: "created_at" },
      { text: "Actions", value: "action", sortable: false },
    ],
    campaigns: [],
    editedIndex: -1,
    editedItem: {
      name: "",
      status: "",
      type: "",
      budget: 0,
    },
    defaultItem: {
      name: "",
      status: "",
      type: "",
      budget: 0,
    },
    status: ["Active", "Paused"],
    type: [
      "Headline Search",
      "Product Display",
      "Sponsored Product",
      "Sponsored Brands",
    ],
    isLoading: true,
  }),

  computed: {
    formTitle() {
      return this.editedIndex === -1 ? "New Campaign" : "Edit Campaign";
    },
  },

  watch: {
    dialog(val) {
      val || this.close();
    },
  },

  created() {
    this.initialize();
  },

  methods: {
    async initialize() {
      try {
        let { data } = await axios.get("http://localhost:8000/api/campaign");
        this.campaigns = data;
      } catch (err) {
        console.log(err);
      } finally {
        this.isLoading = false;
      }
    },

    editItem(id) {
      this.editedIndex = id;
      let temp = this.campaigns.find((obj) => obj.id === id);
      this.editedItem = Object.assign({}, temp);
      this.dialog = true;
    },

    deleteItem(id) {
      if (confirm("Are you sure you want to delete this item?")) {
        this.campaigns = this.campaigns.filter((el) => el.id !== id);
        axios.delete(`http://localhost:8000/api/deletecampaign/${id}`);
      }
    },

    close() {
      this.dialog = false;
      setTimeout(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
        this.editedIndex = -1;
      }, 300);
    },

    save() {
      if (this.editedIndex > -1) {
        let temp = this.campaigns.find((obj) => obj.id === this.editedIndex);
        Object.assign(temp, this.editedItem);
        this.editedItem.budget = parseInt(this.editedItem.budget);
        axios.put(
          `http://localhost:8000/api/updatecampaign/${this.editedIndex}`,
          JSON.stringify(this.editedItem)
        );
        this.close();
      } else {
        let temp = this.editedItem;
        temp.budget = parseInt(temp.budget);
        let date = new Date();
        temp.created_at = date.toString().slice(0, 24);
        temp.id = Math.floor(Math.random() * 1000000000);
        axios.post(
          "http://localhost:8000/api/newcampaign",
          JSON.stringify(temp)
        );
        this.campaigns = [...this.campaigns, temp];
        this.close();
      }
    },

    getColor(status) {
      if (status === "Paused") return "red";
      return "green";
    },
  },
};
</script>
