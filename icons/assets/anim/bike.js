import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const bike = el.querySelectorAll("path, circle");
  animate(bike, { translateX: [0, 3, -1, 0] }, { duration: τmax });

  const head = el.querySelector('circle[cx="15"]');
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
