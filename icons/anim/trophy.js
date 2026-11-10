import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const cup = el.querySelector("g");
  animate(cup, { translateY: [0, -3, 1, 0] }, { duration: τmax });
  const base = el.querySelector('path[d="M4 22h16"]');
  animate(base, { scaleX: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
