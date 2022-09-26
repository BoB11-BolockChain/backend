<fieldset>
	<label>
	<input type = "radio" name = "Index_Select" value = "Create_Challenges" 
	<?php if ($H_GenSel == "Create_Challenges") {echo "checked";}?>/>
	<span>Create Challenges</span>
	</label>

	<label>
	<input type = "radio" name = "Index_Select" value = "Show_Agents"
	<?php if ($H_GenSel == "Show_Agents") {echo "checked";}?>/>
	<span>Show Attack Agents List</span>
	</label>

	<label>
	<input type = "radio" name = "Index_Select" value = "temp1" disabled 
	<?php if ($H_GenSel == "temp1") {echo "checked";}?>/>
	<span>temp1</span>
	</label>

	<label>
	<input type = "radio" name = "Index_Select" value = "temp2" disabled
	<?php if ($H_GenSel == "temp2") {echo "checked";}?>/>
	<span>Temp2</span>
	</label>
</fieldset>
