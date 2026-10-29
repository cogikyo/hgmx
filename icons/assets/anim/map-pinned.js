import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ground = el.querySelector('path[d^="M4.474"]');
  animate(ground, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pin = el.querySelector("g");
  animate(pin, { translateY: [-6, 1, 0], scaleY: [1.15, 0.9, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const dot = el.querySelector("circle");
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
