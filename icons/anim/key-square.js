import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const key = el.querySelectorAll("path");
  animate(key, { translateX: [4, -1, 0], translateY: [-4, 1, 0] }, { duration: τmax });

  const slot = el.querySelector('path[d="m14 7 3 3"]');
  animate(slot, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
