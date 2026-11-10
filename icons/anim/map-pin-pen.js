import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const pin = el.querySelector("g");
  animate(pin, { translateY: [-6, 1, 0], scaleY: [1.15, 0.9, 1] }, { duration: τmid });

  const pen = el.querySelector('path[d^="M21.378"]');
  animate(pen, { rotate: [-15, 15, -10, 0], translateX: [0, -3, 0] }, { duration: τmax });
}
