import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const flask = el.querySelector("g");
  animate(flask, { rotate: [0, -10, 8, -4, 0] }, { duration: τmax });

  const liquid = el.querySelector('path[d^="M5 15"]');
  animate(liquid, { rotate: [0, 12, -8, 0] }, { duration: τmid, delay: τmin });
}
