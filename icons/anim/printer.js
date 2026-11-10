import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const sheet = el.querySelector("rect");
  animate(sheet, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: τmax });

  const body = el.querySelector('path[d^="M6 18"]');
  animate(body, { scaleY: [1, 0.9, 1] }, { duration: τmin, delay: τmid });
}
