import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const reticle = el.querySelector("g");
  animate(reticle, { scale: [1.5, 0.9, 1], opacity: [0, 1] }, { duration: τmax });

  const target = el.querySelector('circle[r="3"]');
  animate(target, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
