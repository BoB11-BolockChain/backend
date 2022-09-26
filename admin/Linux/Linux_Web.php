<fieldset>
        <label>
        <input type = "radio" name = "Attack1_Select" value = "File_Upload"
        <?php if ($H_Attack1 == "File_Upload") {echo "checked";}?>/>
        <span>File Upload</span>
        </label>

        <label>
        <input type = "radio" name = "Attack1_Select" value = "File_Download"
        <?php if ($H_Attack1 == "File_Download") {echo "checked";}?>/>
        <span>File Download</span>
        </label>

        <label>
        <input type = "radio" name = "Attack1_Select" value = "XSS"
        <?php if ($H_Attack1 == "XSS") {echo "checked";}?>/>
        <span>XSS</span>
        </label>

        <label>
        <input type = "radio" name = "Attack1_Select" value = "SQL_injection"
        <?php if ($H_Attack1 == "SQL_injection") {echo "checked";}?>/>
        <span>SQL injection</span>
        </label>

        <label>
        <input type = "radio" name = "Attack1_Select" value = "OS_Command_Injection"
        <?php if ($H_Attack1 == "OS_Command_Injection") {echo "checked";}?>/>
        <span>OS Command Injection</span>
        </label>

        <label>
        <input type = "radio" name = "Attack1_Select" value = "PHP_Injection"
        <?php if ($H_Attack1 == "PHP_Injection") {echo "checked";}?>/>
        <span>PHP Injection</span>
        </label>
</fieldset>
