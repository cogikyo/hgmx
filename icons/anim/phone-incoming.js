import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const handset = el.querySelector('path[d^="M13.832"]');
  animate(handset, { rotate: [0, -10, 6, 0] }, { duration: τmid });

  const arrow = el.querySelector("g");
  animate(arrow, { translateX: [3, -1, 0], translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
