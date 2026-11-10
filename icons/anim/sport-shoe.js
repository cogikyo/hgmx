import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const shoe = el.querySelector("g");
  animate(shoe, { rotate: [0, -10, 4, 0] }, { duration: τmax });
  const dash = el.querySelector('path[d="M19 18h3"]');
  animate(dash, { scaleX: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
