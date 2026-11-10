import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const head = el.querySelector("g");
  animate(head, { rotate: [0, -18, 8, 0] }, { duration: τmax });

  const light = el.querySelector('path[d="M7 9h.01"]');
  animate(light, { opacity: [1, 0, 1] }, { duration: τmin, delay: τmid });
}
