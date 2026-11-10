import { animate, τmax } from "../motion.js";

export default function (el) {
  const innerGear = el.querySelectorAll("circle");
  animate(innerGear, { scale: [0.5, 1.25, 1] }, { duration: τmax });

  const outerGear = el.querySelector("path");
  animate(
    outerGear,
    { rotate: [0, 360], scale: [0, 1.15, 1], opacity: [0, 1] },
    { duration: τmax },
  );
}
