import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const capsule = el.querySelector("rect");
  animate(capsule, { scale: [1, 0.85, 1] }, { duration: τmid });

  const wave = el.querySelector('path[d^="M19 10"]');
  animate(wave, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}
