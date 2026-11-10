import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M6 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const prompt = el.querySelector('path[d^="m8 16"]');
  animate(prompt, { translateX: [-3, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const cursor = el.querySelector('path[d="M12 18h4"]');
  animate(cursor, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
