import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const trend = el.querySelector("g");
  animate(trend, { translateX: [0, 3.5, -1.5, 0], translateY: [0, 3.5, -1.5, 0] }, { duration: τmax });

  const head = el.querySelector('path[d="M16 17h6v-6"]');
  animate(head, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
