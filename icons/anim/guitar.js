import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const guitar = el.querySelector("g");
  animate(guitar, { rotate: [0, -8, 5, 0] }, { duration: τmax });

  const bridge = el.querySelector('path[d="m6 16 2 2"]');
  animate(bridge, { scale: [1, 1.5, 1] }, { duration: τmin, delay: τmid });
}
