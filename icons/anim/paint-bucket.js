import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bucket = el.querySelector("g");
  animate(bucket, { rotate: [0, 12, 0] }, { duration: τmid });

  const drop = el.querySelector('path[d^="M21.145"]');
  animate(drop, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
