import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const left = el.querySelector('path[d^="M9 17"]');
  animate(left, { translateX: [-3, 1, 0] }, { duration: τmid });

  const right = el.querySelector('path[d^="M15 7"]');
  animate(right, { translateX: [3, -1, 0] }, { duration: τmid });

  const bridge = el.querySelector("line");
  animate(bridge, { scaleX: [0, 1] }, { duration: τmid, delay: τmin });
}
