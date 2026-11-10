import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const phone = el.querySelector("rect");
  animate(phone, { translateX: [0, -1, 1, -1, 0] }, { duration: τmid });

  const buzz = el.querySelectorAll("path");
  animate(buzz, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}
