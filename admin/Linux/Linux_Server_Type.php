<fieldset>
        <label>
        <input type = "radio" name = "Server_Type" value = "Web_Server"
        <?php if ($H_ServerType == "Web_Server") {echo "checked";}?>/>
        <span>apache Web Server</span>
        </label>

        <label>
        <input type = "radio" name = "Server_Type" value = "DB_Server"
        <?php if ($H_ServerType == "DB_Server") {echo "checked";}?>/>
        <span>DB Server</span>
        </label>

        <label>
        <input type = "radio" name = "Server_Type" value = "Normal_Server"
        <?php if ($H_ServerType == "Normal_Server") {echo "checked";}?>/>
        <span>Normal Server</span>
        </label>

        <label>
        <input type = "radio" name = "Server_Type" value = "File_Server"
        <?php if ($H_ServerType == "File_Server") {echo "checked";}?>/>
        <span>File Server</span>
        </label>
</fieldset>