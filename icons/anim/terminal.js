import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const prompt = el.querySelector('path[d^="m4 17"]');
  animate(prompt, { translateX: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid });

  const cursor = el.querySelector('path[d^="M12 19"]');
  animate(cursor, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
