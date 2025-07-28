import { PUBLIC_API_URL } from "$env/static/public";
import { redirect } from "@sveltejs/kit";

async function getNewTierlistURL() : Promise<string> {
	const uuid = crypto.randomUUID();

	const response = await fetch(`${PUBLIC_API_URL}/tierlist`, {
		method: 'POST',
		body: JSON.stringify({ 'uuid': uuid, 'name': 'New Tierlist' })
	})

	return `${PUBLIC_WWW_URL}/${uuid}`;
}

export async function load() {
	redirect(307, await getNewTierlistURL());
}
