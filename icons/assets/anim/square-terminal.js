import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const prompt = el.querySelector('path[d^="m7 11"]');
  animate(prompt, { translateX: [-3, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const cursor = el.querySelector('path[d^="M11 13"]');
  animate(cursor, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
