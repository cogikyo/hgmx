import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const corners = el.querySelector("g");
  animate(corners, { scale: [1.25, 0.9, 1] }, { duration: τmax });

  const subject = el.querySelector("circle");
  animate(subject, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
