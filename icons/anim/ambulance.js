import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const ambulance = el.querySelectorAll("path, circle");
  animate(ambulance, { translateX: [-6, 1, 0], opacity: [0, 1] }, { duration: τmax });

  const cross = el.querySelectorAll('path[d^="M10 10"], path[d^="M8 8"]');
  animate(cross, { scale: [1, 0.85, 1.15, 1] }, { duration: τmid, delay: τmin });
}
