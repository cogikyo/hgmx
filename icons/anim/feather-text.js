import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const feather = el.querySelector("g");
  animate(feather, { rotate: [0, -12, 8, 0] }, { duration: τmax });

  const line = el.querySelector('path[d="M2 19h8"]');
  animate(line, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
