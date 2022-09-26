<fieldset>
        <label>
        <input type = "radio" name = "Attack3_Select" value = "read_passwd" 
        <?php if ($H_Attack3 == "read_passwd") {echo "checked";}?>/>
        <span>read /etc/passwd file</span>
        </label>

        <label>
        <input type = "radio" name = "Attack3_Select" value = "Attack2" 
        <?php if ($H_Attack3 == "Attack2") {echo "checked";}?>/>
        <span>Attack2</span>
        </label>

        <label>
        <input type = "radio" name = "Attack3_Select" value = "Attack3" 
        <?php if ($H_Attack3 == "Attack3") {echo "checked";}?>/>
        <span>Attack3</span>
        </label>

        <label>
        <input type = "radio" name = "Attack3_Select" value = "Attack4" 
        <?php if ($H_Attack3 == "Attack4") {echo "checked";}?>/>
        <span>Attack4</span>
        </label>
</fieldset>

