import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const head = el.querySelectorAll("path");
  animate(head, { translateY: [0, 1.5, -1, 0] }, { duration: τmax });

  const eyes = el.querySelectorAll('path[d$="v-.5"]');
  animate(eyes, { scaleY: [1, 0.1, 1] }, { duration: τmin, delay: τmid });
}
