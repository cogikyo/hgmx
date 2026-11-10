import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ground = el.querySelector('path[d^="M5 21"]');
  animate(ground, { scaleX: [0, 1] }, { duration: τmid });

  const sprout = el.querySelector("g");
  animate(sprout, { scale: [0, 1.15, 1] }, { duration: τmid, delay: τmin });
}
